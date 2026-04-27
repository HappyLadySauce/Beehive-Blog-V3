package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/testkit"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/pb"
)

func TestTagDeleteFailsWhenBound(t *testing.T) {
	t.Parallel()

	manager := contentservice.NewManager(testkit.NewServiceDependencies(t))
	actor := contentservice.Actor{UserID: 1, SessionID: 10, Role: "admin"}
	tag, err := manager.CreateTag.Execute(context.Background(), actor, &pb.CreateTagRequest{Name: "Bound", Slug: "bound"})
	if err != nil {
		t.Fatalf("create tag failed: %v", err)
	}
	if _, err := manager.CreateContent.Execute(context.Background(), actor, &pb.CreateContentRequest{
		Type:   pb.ContentType_CONTENT_TYPE_ARTICLE,
		Title:  "Tagged",
		Slug:   "tagged",
		TagIds: []string{tag.Tag.TagId},
	}); err != nil {
		t.Fatalf("create content failed: %v", err)
	}
	_, err = manager.DeleteTag.Execute(context.Background(), actor, &pb.DeleteTagRequest{TagId: tag.Tag.TagId})
	if !errors.Is(err, errs.E(errs.CodeContentTagInUse)) {
		t.Fatalf("expected tag in use, got %v", err)
	}
}
