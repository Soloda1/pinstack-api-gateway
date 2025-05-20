package models

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "github.com/soloda1/pinstack-proto-definitions/gen/go/pinstack-proto-definitions/post/v1"
)

type CreatePostDTO struct {
	AuthorID   int64             `json:"author_id"`
	Title      string            `json:"title"`
	Content    *string           `json:"content,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	MediaItems []*PostMediaInput `json:"media_items,omitempty"`
}

type PostMediaInput struct {
	URL      string    `json:"url"`
	Type     MediaType `json:"type"`
	Position int32     `json:"position"`
}

type Post struct {
	ID        int64     `json:"id"`
	AuthorID  int64     `json:"author_id"`
	Title     string    `json:"title"`
	Content   *string   `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostDetailed struct {
	Post   *Post        `json:"post,omitempty"`
	Author *User        `json:"author,omitempty"`
	Media  []*PostMedia `json:"media,omitempty"`
	Tags   []*Tag       `json:"tags,omitempty"`
}

type PostFilters struct {
	AuthorID      *int64
	TagNames      []string
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	Limit         *int
	Offset        *int
}

type PostMedia struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	URL       string    `json:"url"`
	Type      MediaType `json:"type"`
	Position  int32     `json:"position"`
	CreatedAt time.Time `json:"created_at"`
}

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type PostTag struct {
	PostID int64 `json:"post_id"`
	TagID  int64 `json:"tag_id"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdatePostDTO struct {
	Title      *string           `json:"title,omitempty"`
	Content    *string           `json:"content,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	MediaItems []*PostMediaInput `json:"media_items,omitempty"`
}

func CreatePostDTOToProto(dto *CreatePostDTO) *pb.CreatePostRequest {
	if dto == nil {
		return nil
	}
	media := make([]*pb.MediaInput, 0, len(dto.MediaItems))
	for _, m := range dto.MediaItems {
		media = append(media, &pb.MediaInput{
			Url:      m.URL,
			Type:     string(m.Type),
			Position: m.Position,
		})
	}
	return &pb.CreatePostRequest{
		AuthorId: dto.AuthorID,
		Title:    dto.Title,
		Content:  derefString(dto.Content),
		Tags:     dto.Tags,
		Media:    media,
	}
}

func UpdatePostDTOToProto(id int64, dto *UpdatePostDTO) *pb.UpdatePostRequest {
	if dto == nil {
		return nil
	}
	media := make([]*pb.MediaInput, 0, len(dto.MediaItems))
	for _, m := range dto.MediaItems {
		media = append(media, &pb.MediaInput{
			Url:      m.URL,
			Type:     string(m.Type),
			Position: m.Position,
		})
	}
	return &pb.UpdatePostRequest{
		Id:      id,
		Title:   derefString(dto.Title),
		Content: derefString(dto.Content),
		Tags:    dto.Tags,
		Media:   media,
	}
}

func PostFiltersToProto(filters *PostFilters) *pb.ListPostsRequest {
	if filters == nil {
		return &pb.ListPostsRequest{}
	}
	var authorId int64
	if filters.AuthorID != nil {
		authorId = *filters.AuthorID
	}
	var offset, limit int32
	if filters.Offset != nil {
		offset = int32(*filters.Offset)
	}
	if filters.Limit != nil {
		limit = int32(*filters.Limit)
	}
	return &pb.ListPostsRequest{
		AuthorId: authorId,
		Offset:   offset,
		Limit:    limit,
	}
}

func PostDetailedFromProto(p *pb.Post) *PostDetailed {
	if p == nil {
		return nil
	}
	media := make([]*PostMedia, 0, len(p.Media))
	for _, m := range p.Media {
		media = append(media, &PostMedia{
			ID:        m.Id,
			URL:       m.Url,
			Type:      MediaType(m.Type),
			Position:  m.Position,
			CreatedAt: pbTimestampToTime(m.CreatedAt),
		})
	}
	return &PostDetailed{
		Post: &Post{
			ID:        p.Id,
			AuthorID:  p.AuthorId,
			Title:     p.Title,
			Content:   ptrString(p.Content),
			CreatedAt: pbTimestampToTime(p.CreatedAt),
			UpdatedAt: pbTimestampToTime(p.UpdatedAt),
		},
		Media: media,
		Tags:  tagsFromStrings(p.Tags),
	}
}

func pbTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return time.Unix(ts.Seconds, int64(ts.Nanos)).UTC()
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func tagsFromStrings(tags []string) []*Tag {
	result := make([]*Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, &Tag{Name: t})
	}
	return result
}
