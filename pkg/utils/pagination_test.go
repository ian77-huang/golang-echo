package utils

import "testing"

func TestGetPagination(t *testing.T) {
	for _, tt := range []struct {
		name                    string
		currentPage, totalPages int
		wantValues              []int
		wantCurrent             bool
	}{
		{"single page", 1, 1, []int{1}, true},
		{"middle window", 5, 20, []int{3, 4, 5, 6, 7}, true},
		{"at beginning", 1, 20, []int{1, 2, 3, 4, 5}, true},
		{"at end", 20, 20, []int{16, 17, 18, 19, 20}, true},
		{"fewer pages than spacing", 3, 3, []int{1, 2, 3}, true},
		{"zero total pages", 1, 0, nil, false},
		{"negative total pages", 1, -1, nil, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := GetPagination(tt.currentPage, tt.totalPages, 5)
			wantTotal := tt.totalPages
			if wantTotal <= 0 {
				wantTotal = 0
			}
			if p.CurrentPage != tt.currentPage || p.TotalPages != wantTotal {
				t.Fatalf("unexpected header: %#v", p)
			}
			if len(p.Pages) != len(tt.wantValues) {
				t.Fatalf("pages = %#v, want %#v", p.Pages, tt.wantValues)
			}
			for i, want := range tt.wantValues {
				if p.Pages[i].Value != want {
					t.Fatalf("page %d = %d, want %d", i, p.Pages[i].Value, want)
				}
			}
			current := false
			for _, page := range p.Pages {
				if page.IsCurrentPage {
					current = true
				}
			}
			if current != tt.wantCurrent {
				t.Fatalf("current flag mismatch: got %v want %v", current, tt.wantCurrent)
			}
		})
	}
}

func TestGetPaginationFallsBackToDefaultSpacing(t *testing.T) {
	p := GetPagination(10, 100, 0)
	if len(p.Pages) != 5 {
		t.Fatalf("expected 5 pages with default spacing, got %#v", p.Pages)
	}
	p = GetPagination(10, 100, -3)
	if len(p.Pages) != 5 {
		t.Fatalf("expected 5 pages with negative spacing fallback, got %#v", p.Pages)
	}
}

func TestGetPaginationWithCustomSpacing(t *testing.T) {
	p := GetPagination(10, 50, 3)
	if len(p.Pages) != 3 || p.Pages[0].Value != 9 || p.Pages[2].Value != 11 {
		t.Fatalf("unexpected pages: %#v", p.Pages)
	}
}
