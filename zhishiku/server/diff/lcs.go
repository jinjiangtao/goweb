package diff

type DiffLine struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	OldLine *int   `json:"old_line,omitempty"`
	NewLine *int   `json:"new_line,omitempty"`
}

type DiffResult struct {
	TitleDiff *TextDiff     `json:"title_diff,omitempty"`
	CatDiff   *CategoryDiff `json:"category_diff,omitempty"`
	Lines     []DiffLine    `json:"lines"`
	HasChange bool          `json:"has_change"`
}

type TextDiff struct {
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
	Changed  bool   `json:"changed"`
}

type CategoryDiff struct {
	OldID   uint   `json:"old_id"`
	OldName string `json:"old_name"`
	NewID   uint   `json:"new_id"`
	NewName string `json:"new_name"`
	Changed bool   `json:"changed"`
}

func lcsLength(a, b []string) [][]int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] >= dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp
}

func backtrack(dp [][]int, a, b []string, i, j int) []DiffLine {
	var result []DiffLine
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && a[i-1] == b[j-1] {
			ol := i
			nl := j
			result = append(result, DiffLine{
				Type:    "unchanged",
				Content: a[i-1],
				OldLine: &ol,
				NewLine: &nl,
			})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			nl := j
			result = append(result, DiffLine{
				Type:    "added",
				Content: b[j-1],
				NewLine: &nl,
			})
			j--
		} else {
			ol := i
			result = append(result, DiffLine{
				Type:    "removed",
				Content: a[i-1],
				OldLine: &ol,
			})
			i--
		}
	}

	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}
	return result
}

func ComputeDiff(oldText, newText string) []DiffLine {
	oldLines := splitLines(oldText)
	newLines := splitLines(newText)

	dp := lcsLength(oldLines, newLines)
	return backtrack(dp, oldLines, newLines, len(oldLines), len(newLines))
}

func splitLines(text string) []string {
	if text == "" {
		return nil
	}
	var lines []string
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '\n' {
			lines = append(lines, text[start:i])
			start = i + 1
		}
	}
	if start <= len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}

func HasContentChange(oldText, newText string) bool {
	return oldText != newText
}
