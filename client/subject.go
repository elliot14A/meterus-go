package client

import (
	"context"

	subject "github.com/elliot14A/meterus-go/subject/v1"
)

type SubjectService struct {
	client subject.SubjectServiceClient
	apiKey string
}

func (c *Client) NewSubjectService() *SubjectService {
	return &SubjectService{
		client: subject.NewSubjectServiceClient(c.conn),
		apiKey: c.apiKey,
	}
}

func (s *SubjectService) Create(ctx context.Context, id string, displayName *string) (*subject.Subject, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, s.apiKey)
	return s.client.CreateSubject(ctx, &subject.Subject{
		Id:          id,
		DisplayName: displayName,
	})
}

func (s *SubjectService) GetById(ctx context.Context, id string) (*subject.Subject, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, s.apiKey)
	return s.client.GetSubject(ctx, &subject.SubjectId{SubjectId: id})
}

func (s *SubjectService) ListById(ctx context.Context, page, limit int32) ([]*subject.Subject, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, s.apiKey)
	subjects, err := s.client.ListSubjects(ctx, &subject.ListSubjectRequest{Limit: limit, Page: page})
	if err != nil {
		return nil, err
	}
	return subjects.Subjects, nil
}

func (s *SubjectService) Update(ctx context.Context, id string, displayName *string) (*subject.Subject, error) {
	ctx = AddApiKeyAuthorizationHeader(ctx, s.apiKey)
	return s.client.UpdateSubject(ctx, &subject.Subject{Id: id, DisplayName: displayName})
}

func (s *SubjectService) Delete(ctx context.Context, id string) error {
	ctx = AddApiKeyAuthorizationHeader(ctx, s.apiKey)
	_, err := s.client.DeleteSubject(ctx, &subject.SubjectId{SubjectId: id})
	return err
}
