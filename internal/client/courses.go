// Copyright iximiuz Labs 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetCourse gets a course by name.
func (c *Client) GetCourse(ctx context.Context, name string) (*Course, error) {
	var course Course
	if err := c.GetJSON(ctx, "/courses/"+url.PathEscape(name), nil, &course); err != nil {
		return nil, fmt.Errorf("getting course %q: %w", name, err)
	}
	return &course, nil
}

// ListCourses lists all courses.
func (c *Client) ListCourses(ctx context.Context) ([]Course, error) {
	var courses []Course
	if err := c.GetJSON(ctx, "/courses", nil, &courses); err != nil {
		return nil, fmt.Errorf("listing courses: %w", err)
	}
	return courses, nil
}

// CreateCourse creates a new course.
func (c *Client) CreateCourse(ctx context.Context, req CreateCourseRequest) (*Course, error) {
	var course Course
	if err := c.PostJSON(ctx, "/courses", req, &course); err != nil {
		return nil, fmt.Errorf("creating course %q: %w", req.Name, err)
	}
	return &course, nil
}

// DeleteCourse deletes a course by name.
func (c *Client) DeleteCourse(ctx context.Context, name string) error {
	if err := c.Delete(ctx, "/courses/"+url.PathEscape(name), nil); err != nil {
		return fmt.Errorf("deleting course %q: %w", name, err)
	}
	return nil
}
