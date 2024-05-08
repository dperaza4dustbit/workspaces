package workspace

import (
	"context"
	"fmt"

	workspacesv1alpha1 "github.com/konflux-workspaces/workspaces/operator/api/v1alpha1"
	ccontext "github.com/konflux-workspaces/workspaces/server/core/context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -destination=mocks/mocks.go -package=mocks . WorkspaceCreator

// CreateWorkspaceCommand contains the information needed to create a new workspace
type CreateWorkspaceCommand struct {
	Workspace workspacesv1alpha1.InternalWorkspace
}

// CreateWorkspaceResponse contains the newly-created workspace
type CreateWorkspaceResponse struct {
	Workspace *workspacesv1alpha1.InternalWorkspace
}

type WorkspaceCreator interface {
	CreateUserWorkspace(ctx context.Context, user string, workspace *workspacesv1alpha1.InternalWorkspace, opts ...client.CreateOption) error
}

type CreateWorkspaceHandler struct {
	creator WorkspaceCreator
}

func NewCreateWorkspaceHandler(creator WorkspaceCreator) *CreateWorkspaceHandler {
	return &CreateWorkspaceHandler{creator: creator}
}

func (h *CreateWorkspaceHandler) Handle(ctx context.Context, request CreateWorkspaceCommand) (*CreateWorkspaceResponse, error) {
	u, ok := ctx.Value(ccontext.UserKey).(string)
	if !ok {
		return nil, fmt.Errorf("unauthenticated request")
	}

	// TODO: validate the workspace; maybe punt to a webhook down the line?

	// write the workspace
	workspace := request.Workspace.DeepCopy()
	opts := &client.CreateOptions{}
	if err := h.creator.CreateUserWorkspace(ctx, u, workspace, opts); err != nil {
		return nil, err
	}

	response := &CreateWorkspaceResponse{
		Workspace: workspace,
	}
	return response, nil
}
