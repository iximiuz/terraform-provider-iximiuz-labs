// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

// Playground represents a playground definition.
type Playground struct {
	ID             string                   `json:"id"`
	Owner          string                   `json:"owner"`
	Name           string                   `json:"name"`
	Base           string                   `json:"base,omitempty"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	Categories     []string                 `json:"categories,omitempty"`
	Cover          string                   `json:"cover,omitempty"`
	Markdown       string                   `json:"markdown,omitempty"`
	Published      bool                     `json:"published"`
	PageURL        string                   `json:"pageUrl"`
	Networks       []PlaygroundNetwork      `json:"networks,omitempty"`
	Machines       []PlaygroundMachine      `json:"machines,omitempty"`
	Tabs           []PlaygroundTab          `json:"tabs,omitempty"`
	InitTasks      map[string]InitTask      `json:"initTasks,omitempty"`
	InitConditions *InitConditions          `json:"initConditions,omitempty"`
	RegistryAuth   string                   `json:"registryAuth,omitempty"`
	PortForwards   []PortForward            `json:"portForwards,omitempty"`
	AccessControl  *PlaygroundAccessControl `json:"accessControl,omitempty"`
	UserAccess     *PlaygroundUserAccess    `json:"userAccess,omitempty"`
}

// CreatePlaygroundRequest is the request body for creating a playground.
type CreatePlaygroundRequest struct {
	Name           string                   `json:"name"`
	Base           string                   `json:"base,omitempty"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	Categories     []string                 `json:"categories,omitempty"`
	Markdown       string                   `json:"markdown,omitempty"`
	Networks       []PlaygroundNetwork      `json:"networks,omitempty"`
	Machines       []PlaygroundMachine      `json:"machines,omitempty"`
	Tabs           []PlaygroundTab          `json:"tabs,omitempty"`
	InitTasks      map[string]InitTask      `json:"initTasks,omitempty"`
	InitConditions *InitConditions          `json:"initConditions,omitempty"`
	RegistryAuth   string                   `json:"registryAuth,omitempty"`
	PortForwards   []PortForward            `json:"portForwards,omitempty"`
	AccessControl  *PlaygroundAccessControl `json:"accessControl,omitempty"`
}

// UpdatePlaygroundRequest is the request body for updating a playground.
// Name and Base are immutable and not included.
type UpdatePlaygroundRequest struct {
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
	Categories     []string                 `json:"categories,omitempty"`
	Cover          string                   `json:"cover,omitempty"`
	Markdown       string                   `json:"markdown,omitempty"`
	Networks       []PlaygroundNetwork      `json:"networks,omitempty"`
	Machines       []PlaygroundMachine      `json:"machines,omitempty"`
	Tabs           []PlaygroundTab          `json:"tabs,omitempty"`
	InitTasks      map[string]InitTask      `json:"initTasks,omitempty"`
	InitConditions *InitConditions          `json:"initConditions,omitempty"`
	RegistryAuth   string                   `json:"registryAuth,omitempty"`
	PortForwards   []PortForward            `json:"portForwards,omitempty"`
	AccessControl  *PlaygroundAccessControl `json:"accessControl,omitempty"`
}

// PlaygroundNetwork defines a network in a playground.
type PlaygroundNetwork struct {
	Name    string `json:"name"`
	Subnet  string `json:"subnet,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Private bool   `json:"private,omitempty"`
}

// PlaygroundMachine defines a machine in a playground.
type PlaygroundMachine struct {
	Name         string               `json:"name"`
	Users        []MachineUser        `json:"users,omitempty"`
	Kernel       *MachineKernel       `json:"kernel,omitempty"`
	Drives       []MachineDrive       `json:"drives,omitempty"`
	Network      *MachineNetwork      `json:"network,omitempty"`
	Resources    *MachineResources    `json:"resources,omitempty"`
	StartupFiles []MachineStartupFile `json:"startupFiles,omitempty"`
	NoSSH        bool                 `json:"noSsh,omitempty"`
}

// MachineUser defines a user on a machine.
type MachineUser struct {
	Name    string `json:"name"`
	Default bool   `json:"default,omitempty"`
	Welcome string `json:"welcome,omitempty"`
}

// MachineKernel defines the kernel configuration for a machine.
type MachineKernel struct {
	Source   string          `json:"source,omitempty"`
	Snapshot *RemoteSnapshot `json:"snapshot,omitempty"`
}

// MachineDrive defines a drive attached to a machine.
type MachineDrive struct {
	Source     string          `json:"source,omitempty"`
	Mount      string          `json:"mount,omitempty"`
	Size       string          `json:"size,omitempty"`
	Filesystem string          `json:"filesystem,omitempty"`
	ReadOnly   bool            `json:"readOnly,omitempty"`
	Persistent bool            `json:"persistent,omitempty"`
	Snapshot   *RemoteSnapshot `json:"snapshot,omitempty"`
}

// RemoteSnapshot represents a remote snapshot reference.
type RemoteSnapshot struct {
	ID      string `json:"id"`
	LeaseID string `json:"leaseId,omitempty"`
}

// MachineNetwork defines the network configuration for a machine.
type MachineNetwork struct {
	Interfaces []MachineNetworkInterface `json:"interfaces,omitempty"`
}

// MachineNetworkInterface defines a network interface on a machine.
type MachineNetworkInterface struct {
	Address string `json:"address"`
	Network string `json:"network"`
}

// MachineResources defines resource limits for a machine.
type MachineResources struct {
	CPUCount int    `json:"cpuCount,omitempty"`
	RAMSize  string `json:"ramSize,omitempty"`
}

// MachineStartupFile defines a file to be created at machine startup.
type MachineStartupFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Mode    string `json:"mode,omitempty"`
	Owner   string `json:"owner,omitempty"`
	Append  bool   `json:"append,omitempty"`
}

// PlaygroundTab defines a tab in the playground UI.
type PlaygroundTab struct {
	ID          string `json:"id,omitempty"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Machine     string `json:"machine,omitempty"`
	Number      int    `json:"number,omitempty"`
	Access      string `json:"access,omitempty"`
	TLS         bool   `json:"tls,omitempty"`
	HostRewrite string `json:"hostRewrite,omitempty"`
	PathRewrite string `json:"pathRewrite,omitempty"`
	URL         string `json:"url,omitempty"`
}

// InitTask defines an initialization task.
type InitTask struct {
	Name           string          `json:"name,omitempty"`
	Machine        string          `json:"machine,omitempty"`
	Init           bool            `json:"init,omitempty"`
	User           string          `json:"user,omitempty"`
	TimeoutSeconds int             `json:"timeoutSeconds,omitempty"`
	Needs          []string        `json:"needs,omitempty"`
	Run            string          `json:"run,omitempty"`
	Status         int             `json:"status,omitempty"`
	Conditions     []InitCondition `json:"conditions,omitempty"`
}

// InitCondition defines a condition for an init task.
type InitCondition struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// InitConditions defines initialization condition values.
type InitConditions struct {
	Values []InitConditionValue `json:"values,omitempty"`
}

// InitConditionValue defines a configurable init condition.
type InitConditionValue struct {
	Key      string   `json:"key"`
	Default  string   `json:"default,omitempty"`
	Nullable bool     `json:"nullable,omitempty"`
	Options  []string `json:"options,omitempty"`
}

// PortForward defines a port forward configuration.
type PortForward struct {
	Kind       string `json:"kind"`
	Machine    string `json:"machine"`
	LocalHost  string `json:"localHost,omitempty"`
	LocalPort  int    `json:"localPort,omitempty"`
	RemotePort int    `json:"remotePort,omitempty"`
	RemoteHost string `json:"remoteHost,omitempty"`
}

// PlaygroundAccessControl defines access control for a playground.
type PlaygroundAccessControl struct {
	CanList  []string `json:"canList,omitempty"`
	CanRead  []string `json:"canRead,omitempty"`
	CanStart []string `json:"canStart,omitempty"`
}

// PlaygroundUserAccess represents the current user's access to a playground (read-only).
type PlaygroundUserAccess struct {
	CanList  bool `json:"canList"`
	CanRead  bool `json:"canRead"`
	CanStart bool `json:"canStart"`
}

// Play represents a running instance of a playground.
type Play struct {
	ID          string              `json:"id"`
	CreatedAt   string              `json:"createdAt"`
	UpdatedAt   string              `json:"updatedAt"`
	LastStateAt string              `json:"lastStateAt,omitempty"`
	ExpiresIn   int                 `json:"expiresIn"`
	Status      *PlayStatus         `json:"status,omitempty"`
	Playground  Playground          `json:"playground"`
	Machines    []Machine           `json:"machines,omitempty"`
	Tasks       map[string]PlayTask `json:"tasks,omitempty"`
	PageURL     string              `json:"pageUrl"`

	TutorialName  string `json:"tutorialName,omitempty"`
	ChallengeName string `json:"challengeName,omitempty"`
	CourseName    string `json:"courseName,omitempty"`
	LessonPath    string `json:"lessonPath,omitempty"`
}

// PlayStateEvent represents a single state transition event in a play's lifecycle.
type PlayStateEvent struct {
	State string `json:"state"`
	At    string `json:"at"`
}

// PlayStatusCondition represents a condition in the play status.
type PlayStatusCondition struct {
	Name             string `json:"name"`
	Status           string `json:"status"`
	LastTransitionAt string `json:"lastTransitionAt"`
}

// PlayMachineStatus represents the status of a machine within a play.
type PlayMachineStatus struct {
	Name       string                `json:"name"`
	State      string                `json:"state"`
	Conditions []PlayStatusCondition `json:"conditions,omitempty"`
}

// PlayStatus represents the status of a play as returned by the API.
// The current state is derived from the last element of StateEvents.
type PlayStatus struct {
	StateEvents []PlayStateEvent      `json:"stateEvents,omitempty"`
	Conditions  []PlayStatusCondition `json:"conditions,omitempty"`
	Machines    []PlayMachineStatus   `json:"machines,omitempty"`
	FactoryID   string                `json:"factoryId,omitempty"`
}

// CurrentState returns the current state of the play by looking at the last
// state event. Returns an empty string if there are no state events.
func (s *PlayStatus) CurrentState() string {
	if s == nil || len(s.StateEvents) == 0 {
		return ""
	}
	return s.StateEvents[len(s.StateEvents)-1].State
}

// Machine represents a running machine in a play.
type Machine struct {
	Name      string            `json:"name"`
	Users     []MachineUser     `json:"users,omitempty"`
	Resources *MachineResources `json:"resources,omitempty"`
}

// PlayTask represents the status of a task in a play.
type PlayTask struct {
	Name    string `json:"name"`
	Init    bool   `json:"init,omitempty"`
	Helper  bool   `json:"helper,omitempty"`
	Status  int    `json:"status"`
	Version int    `json:"version,omitempty"`
}

// CreatePlayRequest is the request body for creating a play.
type CreatePlayRequest struct {
	Playground              string              `json:"playground"`
	Tabs                    []PlaygroundTab     `json:"tabs,omitempty"`
	Networks                []PlaygroundNetwork `json:"networks,omitempty"`
	Machines                []PlaygroundMachine `json:"machines,omitempty"`
	InitTasks               map[string]InitTask `json:"initTasks,omitempty"`
	InitConditions          map[string]string   `json:"initConditions,omitempty"`
	SafetyDisclaimerConsent bool                `json:"safetyDisclaimerConsent,omitempty"`
	AsFreeTierUser          bool                `json:"asFreeTierUser,omitempty"`
}

// PlayActionRequest is the request body for play lifecycle actions.
type PlayActionRequest struct {
	Action string `json:"action"`
}

// Play states.
const (
	PlayStateCreated    = "CREATED"
	PlayStateWarmingUp  = "WARMING_UP"
	PlayStateWarmedUp   = "WARMED_UP"
	PlayStateStarting   = "STARTING"
	PlayStateRunning    = "RUNNING"
	PlayStateStopping   = "STOPPING"
	PlayStateStopped    = "STOPPED"
	PlayStateDestroying = "DESTROYING"
	PlayStateDestroyed  = "DESTROYED"
	PlayStateFailed     = "FAILED"
)

// Challenge represents a challenge content resource.
type Challenge struct {
	CreatedAt       string              `json:"createdAt"`
	UpdatedAt       string              `json:"updatedAt"`
	Name            string              `json:"name"`
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	Categories      []string            `json:"categories,omitempty"`
	Tags            []string            `json:"tags,omitempty"`
	Authors         []Author            `json:"authors,omitempty"`
	PageURL         string              `json:"pageUrl"`
	AttemptCount    int                 `json:"attemptCount"`
	CompletionCount int                 `json:"completionCount"`
	Play            *Play               `json:"play,omitempty"`
	Tasks           map[string]PlayTask `json:"tasks,omitempty"`
}

// Tutorial represents a tutorial content resource.
type Tutorial struct {
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Categories      []string `json:"categories,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Authors         []Author `json:"authors,omitempty"`
	PageURL         string   `json:"pageUrl"`
	AttemptCount    int      `json:"attemptCount"`
	CompletionCount int      `json:"completionCount"`
	Play            *Play    `json:"play,omitempty"`
}

// Course represents a course content resource.
type Course struct {
	CreatedAt string         `json:"createdAt"`
	UpdatedAt string         `json:"updatedAt"`
	Name      string         `json:"name"`
	Title     string         `json:"title"`
	PageURL   string         `json:"pageUrl"`
	Authors   []Author       `json:"authors,omitempty"`
	Modules   []CourseModule `json:"modules,omitempty"`
}

// CourseModule represents a module within a course.
type CourseModule struct {
	Name    string         `json:"name"`
	Title   string         `json:"title"`
	Slug    string         `json:"slug,omitempty"`
	Path    string         `json:"path,omitempty"`
	Lessons []CourseLesson `json:"lessons,omitempty"`
}

// CourseLesson represents a lesson within a course module.
type CourseLesson struct {
	Name       string              `json:"name"`
	Title      string              `json:"title"`
	Slug       string              `json:"slug,omitempty"`
	Path       string              `json:"path,omitempty"`
	Playground *CoursePlayground   `json:"playground,omitempty"`
	Tasks      map[string]PlayTask `json:"tasks,omitempty"`
}

// CoursePlayground represents a playground reference in a course lesson.
type CoursePlayground struct {
	Name string `json:"name"`
}

// CourseVariant represents the type of course.
type CourseVariant string

const (
	CourseVariantSimple  CourseVariant = "simple"
	CourseVariantModular CourseVariant = "modular"
)

// CreateCourseRequest is the request body for creating a course.
type CreateCourseRequest struct {
	Name    string        `json:"name"`
	Variant CourseVariant `json:"variant,omitempty"`
	Sample  bool          `json:"sample,omitempty"`
}

// Roadmap represents a roadmap content resource.
type Roadmap struct {
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	PageURL   string   `json:"pageUrl"`
	Authors   []Author `json:"authors,omitempty"`
}

// SkillPath represents a skill path content resource.
type SkillPath struct {
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	PageURL   string   `json:"pageUrl"`
	Authors   []Author `json:"authors,omitempty"`
}

// Training represents a training content resource.
type Training struct {
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Name      string   `json:"name"`
	Title     string   `json:"title"`
	PageURL   string   `json:"pageUrl"`
	Authors   []Author `json:"authors,omitempty"`
}

// Author represents an author profile.
type Author struct {
	UserID             string `json:"userId"`
	DisplayName        string `json:"displayName"`
	ExternalProfileURL string `json:"externalProfileUrl,omitempty"`
	Official           bool   `json:"official,omitempty"`
}

// CreateAuthorRequest is the request body for creating an author.
type CreateAuthorRequest struct {
	DisplayName        string `json:"displayName"`
	ExternalProfileURL string `json:"externalProfileUrl,omitempty"`
}

// Me represents the current authenticated user.
type Me struct {
	ID              string         `json:"id"`
	GithubProfileID string         `json:"githubProfileId,omitempty"`
	PremiumAccess   *PremiumAccess `json:"premiumAccess,omitempty"`
}

// PremiumAccess represents premium subscription info.
type PremiumAccess struct {
	Until    string `json:"until,omitempty"`
	Lifetime bool   `json:"lifetime,omitempty"`
	Trial    bool   `json:"trial,omitempty"`
}

// Port represents an exposed port on a play machine.
type Port struct {
	ID          string `json:"id"`
	PlayID      string `json:"playId"`
	Machine     string `json:"machine"`
	Number      int    `json:"number"`
	Hostname    string `json:"hostname,omitempty"`
	AccessMode  string `json:"access"`
	TLS         bool   `json:"tls,omitempty"`
	HostRewrite string `json:"hostRewrite,omitempty"`
	PathRewrite string `json:"pathRewrite,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ExposePortRequest is the request body for exposing a port.
type ExposePortRequest struct {
	Machine     string `json:"machine"`
	Number      int    `json:"number"`
	Access      string `json:"access"`
	TLS         bool   `json:"tls,omitempty"`
	HostRewrite string `json:"hostRewrite,omitempty"`
	PathRewrite string `json:"pathRewrite,omitempty"`
}

// Shell represents an exposed shell on a play machine.
type Shell struct {
	ID         string `json:"id"`
	PlayID     string `json:"playId"`
	Machine    string `json:"machine"`
	User       string `json:"user,omitempty"`
	Hostname   string `json:"hostname,omitempty"`
	AccessMode string `json:"access"`
	URL        string `json:"url,omitempty"`
}

// ExposeShellRequest is the request body for exposing a shell.
type ExposeShellRequest struct {
	Machine string `json:"machine"`
	User    string `json:"user,omitempty"`
	Access  string `json:"access"`
}

// CreateContentRequest is a generic request body for creating simple content resources.
type CreateContentRequest struct {
	Name string `json:"name"`
}
