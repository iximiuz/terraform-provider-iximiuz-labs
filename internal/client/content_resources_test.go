// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Challenge Tests ---

func TestGetChallenge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/challenges/my-challenge" {
			t.Errorf("expected path /challenges/my-challenge, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "my-challenge",
			"title": "My Challenge",
			"createdAt": "2025-01-01T00:00:00Z",
			"updatedAt": "2025-01-01T00:00:00Z",
			"pageUrl": "https://labs.iximiuz.com/challenges/my-challenge"
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	ch, err := c.GetChallenge(context.Background(), "my-challenge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch.Name != "my-challenge" {
		t.Errorf("expected name %q, got %q", "my-challenge", ch.Name)
	}
	if ch.Title != "My Challenge" {
		t.Errorf("expected title %q, got %q", "My Challenge", ch.Title)
	}
}

func TestListChallenges(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/challenges" {
			t.Errorf("expected path /challenges, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"ch-1"},{"name":"ch-2"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	challenges, err := c.ListChallenges(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(challenges) != 2 {
		t.Fatalf("expected 2 challenges, got %d", len(challenges))
	}
}

func TestCreateChallenge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var req CreateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-challenge" {
			t.Errorf("expected name %q, got %q", "new-challenge", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-challenge","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	ch, err := c.CreateChallenge(context.Background(), "new-challenge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch.Name != "new-challenge" {
		t.Errorf("expected name %q, got %q", "new-challenge", ch.Name)
	}
}

func TestDeleteChallenge(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/challenges/my-challenge" {
			t.Errorf("expected path /challenges/my-challenge, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteChallenge(context.Background(), "my-challenge")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Tutorial Tests ---

func TestGetTutorial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tutorials/my-tutorial" {
			t.Errorf("expected path /tutorials/my-tutorial, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"my-tutorial","title":"My Tutorial","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":"https://labs.iximiuz.com/tutorials/my-tutorial"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tut, err := c.GetTutorial(context.Background(), "my-tutorial")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tut.Name != "my-tutorial" {
		t.Errorf("expected name %q, got %q", "my-tutorial", tut.Name)
	}
}

func TestListTutorials(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"t-1"},{"name":"t-2"},{"name":"t-3"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tutorials, err := c.ListTutorials(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tutorials) != 3 {
		t.Fatalf("expected 3 tutorials, got %d", len(tutorials))
	}
}

func TestCreateTutorial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-tutorial" {
			t.Errorf("expected name %q, got %q", "new-tutorial", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-tutorial","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tut, err := c.CreateTutorial(context.Background(), "new-tutorial")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tut.Name != "new-tutorial" {
		t.Errorf("expected name %q, got %q", "new-tutorial", tut.Name)
	}
}

func TestDeleteTutorial(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tutorials/my-tutorial" {
			t.Errorf("expected path /tutorials/my-tutorial, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteTutorial(context.Background(), "my-tutorial")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Course Tests ---

func TestGetCourse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/courses/my-course" {
			t.Errorf("expected path /courses/my-course, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"name": "my-course",
			"title": "My Course",
			"createdAt": "2025-01-01T00:00:00Z",
			"updatedAt": "2025-01-01T00:00:00Z",
			"pageUrl": "https://labs.iximiuz.com/courses/my-course",
			"modules": [{"name": "mod-1", "title": "Module 1"}]
		}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	course, err := c.GetCourse(context.Background(), "my-course")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if course.Name != "my-course" {
		t.Errorf("expected name %q, got %q", "my-course", course.Name)
	}
	if len(course.Modules) != 1 {
		t.Fatalf("expected 1 module, got %d", len(course.Modules))
	}
	if course.Modules[0].Name != "mod-1" {
		t.Errorf("expected module name %q, got %q", "mod-1", course.Modules[0].Name)
	}
}

func TestListCourses(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"c-1"},{"name":"c-2"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	courses, err := c.ListCourses(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(courses) != 2 {
		t.Fatalf("expected 2 courses, got %d", len(courses))
	}
}

func TestCreateCourse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateCourseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-course" {
			t.Errorf("expected name %q, got %q", "new-course", req.Name)
		}
		if req.Variant != CourseVariantModular {
			t.Errorf("expected variant %q, got %q", CourseVariantModular, req.Variant)
		}
		if !req.Sample {
			t.Error("expected sample to be true")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-course","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	course, err := c.CreateCourse(context.Background(), CreateCourseRequest{
		Name:    "new-course",
		Variant: CourseVariantModular,
		Sample:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if course.Name != "new-course" {
		t.Errorf("expected name %q, got %q", "new-course", course.Name)
	}
}

func TestDeleteCourse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/courses/my-course" {
			t.Errorf("expected path /courses/my-course, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteCourse(context.Background(), "my-course")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Roadmap Tests ---

func TestGetRoadmap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/roadmaps/my-roadmap" {
			t.Errorf("expected path /roadmaps/my-roadmap, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"my-roadmap","title":"My Roadmap","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":"https://labs.iximiuz.com/roadmaps/my-roadmap"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	r, err := c.GetRoadmap(context.Background(), "my-roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "my-roadmap" {
		t.Errorf("expected name %q, got %q", "my-roadmap", r.Name)
	}
	if r.Title != "My Roadmap" {
		t.Errorf("expected title %q, got %q", "My Roadmap", r.Title)
	}
}

func TestListRoadmaps(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"r-1"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	roadmaps, err := c.ListRoadmaps(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roadmaps) != 1 {
		t.Fatalf("expected 1 roadmap, got %d", len(roadmaps))
	}
}

func TestCreateRoadmap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-roadmap" {
			t.Errorf("expected name %q, got %q", "new-roadmap", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-roadmap","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	rm, err := c.CreateRoadmap(context.Background(), "new-roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rm.Name != "new-roadmap" {
		t.Errorf("expected name %q, got %q", "new-roadmap", rm.Name)
	}
}

func TestDeleteRoadmap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/roadmaps/my-roadmap" {
			t.Errorf("expected path /roadmaps/my-roadmap, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteRoadmap(context.Background(), "my-roadmap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Skill Path Tests ---

func TestGetSkillPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/skill-paths/my-skill-path" {
			t.Errorf("expected path /skill-paths/my-skill-path, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"my-skill-path","title":"My Skill Path","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":"https://labs.iximiuz.com/skill-paths/my-skill-path"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	sp, err := c.GetSkillPath(context.Background(), "my-skill-path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.Name != "my-skill-path" {
		t.Errorf("expected name %q, got %q", "my-skill-path", sp.Name)
	}
	if sp.Title != "My Skill Path" {
		t.Errorf("expected title %q, got %q", "My Skill Path", sp.Title)
	}
}

func TestListSkillPaths(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"sp-1"},{"name":"sp-2"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	paths, err := c.ListSkillPaths(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 skill paths, got %d", len(paths))
	}
}

func TestCreateSkillPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-skill-path" {
			t.Errorf("expected name %q, got %q", "new-skill-path", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-skill-path","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	sp, err := c.CreateSkillPath(context.Background(), "new-skill-path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sp.Name != "new-skill-path" {
		t.Errorf("expected name %q, got %q", "new-skill-path", sp.Name)
	}
}

func TestDeleteSkillPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/skill-paths/my-skill-path" {
			t.Errorf("expected path /skill-paths/my-skill-path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteSkillPath(context.Background(), "my-skill-path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- Training Tests ---

func TestGetTraining(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trainings/my-training" {
			t.Errorf("expected path /trainings/my-training, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name":"my-training","title":"My Training","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":"https://labs.iximiuz.com/trainings/my-training"}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tr, err := c.GetTraining(context.Background(), "my-training")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Name != "my-training" {
		t.Errorf("expected name %q, got %q", "my-training", tr.Name)
	}
}

func TestListTrainings(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"t-1"}]`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	trainings, err := c.ListTrainings(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(trainings) != 1 {
		t.Fatalf("expected 1 training, got %d", len(trainings))
	}
}

func TestCreateTraining(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}
		if req.Name != "new-training" {
			t.Errorf("expected name %q, got %q", "new-training", req.Name)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"name":"new-training","title":"","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","pageUrl":""}`)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	tr, err := c.CreateTraining(context.Background(), "new-training")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Name != "new-training" {
		t.Errorf("expected name %q, got %q", "new-training", tr.Name)
	}
}

func TestDeleteTraining(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/trainings/my-training" {
			t.Errorf("expected path /trainings/my-training, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	err := c.DeleteTraining(context.Background(), "my-training")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
