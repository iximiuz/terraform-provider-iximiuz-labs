// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"path/filepath"

	"github.com/iximiuz/terraform-provider-iximiuz-labs/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"gopkg.in/yaml.v3"
)

// Ensure IximiuzLabsProvider satisfies the provider interface.
var _ provider.Provider = &IximiuzLabsProvider{}

// IximiuzLabsProvider defines the provider implementation.
type IximiuzLabsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// IximiuzLabsProviderModel describes the provider data model.
type IximiuzLabsProviderModel struct {
	APIURL      types.String `tfsdk:"api_url"`
	SessionID   types.String `tfsdk:"session_id"`
	AccessToken types.String `tfsdk:"access_token"`
}

// labctlConfig represents the relevant fields from ~/.iximiuz/labctl/config.yaml.
type labctlConfig struct {
	SessionID   string `yaml:"session_id"`
	AccessToken string `yaml:"access_token"`
}

func (p *IximiuzLabsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "iximiuz-labs"
	resp.Version = p.version
}

func (p *IximiuzLabsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The iximiuz Labs provider allows you to manage resources on the iximiuz Labs platform (https://labs.iximiuz.com).",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The iximiuz Labs API URL. Defaults to `https://labs.iximiuz.com/api`. Can also be set with the `IXIMIUZ_API_URL` environment variable.",
				Optional:            true,
			},
			"session_id": schema.StringAttribute{
				MarkdownDescription: "The session ID for API authentication. Can also be set with the `IXIMIUZ_SESSION_ID` environment variable, or read from the labctl config file.",
				Optional:            true,
				Sensitive:           true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "The access token for API authentication. Can also be set with the `IXIMIUZ_ACCESS_TOKEN` environment variable, or read from the labctl config file.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *IximiuzLabsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data IximiuzLabsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve API URL: config > env > default
	apiURL := "https://labs.iximiuz.com/api"
	if !data.APIURL.IsNull() && !data.APIURL.IsUnknown() {
		apiURL = data.APIURL.ValueString()
	} else if v := os.Getenv("IXIMIUZ_API_URL"); v != "" {
		apiURL = v
	}

	// Resolve session ID: config > env
	sessionID := ""
	if !data.SessionID.IsNull() && !data.SessionID.IsUnknown() {
		sessionID = data.SessionID.ValueString()
	} else if v := os.Getenv("IXIMIUZ_SESSION_ID"); v != "" {
		sessionID = v
	}

	// Resolve access token: config > env
	accessToken := ""
	if !data.AccessToken.IsNull() && !data.AccessToken.IsUnknown() {
		accessToken = data.AccessToken.ValueString()
	} else if v := os.Getenv("IXIMIUZ_ACCESS_TOKEN"); v != "" {
		accessToken = v
	}

	// Fallback to labctl config file if credentials not provided
	if sessionID == "" || accessToken == "" {
		if cfg, err := readLabctlConfig(); err == nil {
			if sessionID == "" {
				sessionID = cfg.SessionID
			}
			if accessToken == "" {
				accessToken = cfg.AccessToken
			}
			if sessionID != "" || accessToken != "" {
				tflog.Debug(ctx, "Using credentials from labctl config file")
			}
		}
	}

	if sessionID == "" || accessToken == "" {
		resp.Diagnostics.AddError(
			"Missing Authentication Credentials",
			"The provider requires both session_id and access_token to authenticate with the iximiuz Labs API. "+
				"Set them in the provider configuration block, via environment variables (IXIMIUZ_SESSION_ID, IXIMIUZ_ACCESS_TOKEN), "+
				"or by running `labctl auth login` to populate the config file at ~/.iximiuz/labctl/config.yaml.",
		)
		return
	}

	// Derive base URL from API URL (strip /api suffix)
	baseURL := apiURL
	if len(apiURL) > 4 && apiURL[len(apiURL)-4:] == "/api" {
		baseURL = apiURL[:len(apiURL)-4]
	}

	c := client.NewClient(client.ClientConfig{
		BaseURL:     baseURL,
		APIBaseURL:  apiURL,
		SessionID:   sessionID,
		AccessToken: accessToken,
		UserAgent:   "terraform-provider-iximiuz-labs/" + p.version,
	})

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *IximiuzLabsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPlaygroundResource,
		NewPlayResource,
		NewChallengeResource,
		NewTutorialResource,
		NewCourseResource,
		NewRoadmapResource,
		NewSkillPathResource,
		NewTrainingResource,
		NewPlayPortResource,
		NewPlayShellResource,
	}
}

func (p *IximiuzLabsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPlaygroundDataSource,
		NewPlaygroundsDataSource,
		NewPlayDataSource,
		NewPlaysDataSource,
		NewMeDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IximiuzLabsProvider{
			version: version,
		}
	}
}

// readLabctlConfig reads credentials from the labctl config file.
func readLabctlConfig() (*labctlConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".iximiuz", "labctl", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg labctlConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
