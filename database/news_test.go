package database

import (
	"testing"
	"time"

	"github.com/pchchv/hcms/models"
)

func TestUpdateNews(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	id, _ := CreateNews(d, newNewsItem("Оригинальный заголовок"))
	item, _ := GetNews(d, int(id))
	item.Title = "Обновлённый заголовок"
	if err := UpdateNews(d, item); err != nil {
		t.Fatalf("UpdateNews: %v", err)
	}

	got, _ := GetNews(d, int(id))
	if got.Title != "Обновлённый заголовок" {
		t.Errorf("expected updated title, got %q", got.Title)
	}
}

func TestCountNews(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	_, _ = CreateNews(d, newNewsItem("A"))
	_, _ = CreateNews(d, newNewsItem("B"))
	count, err := CountNews(d)
	if err != nil {
		t.Fatalf("CountNews: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestGetNews_NotFound(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	got, err := GetNews(d, 9999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Error("expected nil for missing news")
	}
}

func TestDeleteNews_NotFound(t *testing.T) {
	d := openTestDB(t)
	if err := Migrate(d); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	deleted, err := DeleteNews(d, 9999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if deleted != nil {
		t.Error("expected nil for missing news")
	}
}

func newNewsItem(title string) *models.News {
	return &models.News{
		Date:        time.Now().UTC(),
		Title:       title,
		Image:       "",
		Announce:    "Краткое описание " + title,
		Description: "Полное описание " + title,
	}
}
