package main

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// --- CONSTANTS ---
const (
	DEBOUNCE_DURATION = 300 * time.Millisecond
	StateWelcome      = -1
	StateSearch       = 0
	StateBrowse       = 1
	StateEpisodes     = 2
	HistoryFile       = "/Users/huynhngocanhthu/ophim-cli/history.json" //Làm ơn sửa cái path này dùm? Trỏ về cái path của cái project này trên máy cưng nha
)

// --- STYLE ---
var (
	appStyle        = lipgloss.NewStyle().Padding(1, 2)
	leftPaneStyle   = lipgloss.NewStyle().PaddingRight(2).Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(lipgloss.Color("#585b70"))
	rightPaneStyle  = lipgloss.NewStyle().PaddingLeft(2)
	titleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#1e1e2e")).Background(lipgloss.Color("#cba6f7")).Padding(0, 1).Bold(true)
	movieTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#585b70"))
	descStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4"))
	posterStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#f38ba8"))
	welcomeStyle    = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(lipgloss.Color("#b4befe")).Padding(1, 2).Align(lipgloss.Center)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555")).Bold(true).Padding(1)
)

// --- DATA STRUCTURES ---
type History struct {
	LastMovieName string    `json:"last_movie_name"`
	LastMovieSlug string    `json:"last_movie_slug"`
	LastEpName    string    `json:"last_ep_name"`
	LastEpLink    string    `json:"last_ep_link"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Movie struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Year     int    `json:"year"`
	ThumbURL string `json:"thumb_url"`
}

func (m Movie) FilterValue() string { return m.Name }
func (m Movie) Title() string       { return m.Name }
func (m Movie) Description() string { return fmt.Sprintf("%d | %s", m.Year, m.Slug) }

type Episode struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	LinkM3u8  string `json:"link_m3u8"`
	LinkEmbed string `json:"link_embed"`
}

func (e Episode) FilterValue() string { return e.Name }
func (e Episode) Title() string       { return "📺 Tập " + e.Name }
func (e Episode) Description() string { return "OPhim Stream" }

type SearchResponse struct {
	Data struct {
		Items []Movie `json:"items"`
	} `json:"data"`
}

type MovieDetail struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	ASCIIArt string
	Episodes []struct {
		ServerData []Episode `json:"server_data"`
	} `json:"episodes"`
}

type MovieDetailResponse struct {
	Data struct {
		Item MovieDetail `json:"item"`
	} `json:"data"`
}

type DetailLoadedMsg struct {
	Slug    string
	Detail  MovieDetail
	IsError bool
}

// --- MODEL ---
type model struct {
	state        int
	textInput    textinput.Model
	movieList    list.Model
	episodeList  list.Model
	viewport     viewport.Model
	history      History
	selectedSlug string
	err          error

	detailCache   map[string]MovieDetail
	domainCDN     string
	debounceTimer *time.Timer
	width, height int
}

// --- PERSISTENCE ---
func saveHistory(h History) {
	data, _ := json.Marshal(h)
	_ = os.WriteFile(HistoryFile, data, 0644)
}

func loadHistory() History {
	data, err := os.ReadFile(HistoryFile)
	if err != nil {
		return History{}
	}
	var h History
	_ = json.Unmarshal(data, &h)
	return h
}

func initialModel() model {
	_ = godotenv.Load()
	hist := loadHistory()

	ti := textinput.New()
	ti.Placeholder = "Cục dzàng mún coi phim gì nè..."
	ti.Focus()
	ti.Width = 30

	mList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	mList.Title = "Kết quả tìm kiếm"
	mList.SetShowHelp(false)

	eList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	eList.SetShowHelp(false)

	startState := StateSearch
	if hist.LastMovieSlug != "" {
		startState = StateWelcome
	}

	return model{
		state:       startState,
		textInput:   ti,
		movieList:   mList,
		episodeList: eList,
		history:     hist,
		detailCache: make(map[string]MovieDetail),
		domainCDN:   os.Getenv("OPHIM_CDN_IMAGE"),
		viewport:    viewport.New(0, 0),
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		listWidth := int(float64(m.width) * 0.35)
		m.movieList.SetSize(listWidth, m.height-4)
		m.episodeList.SetSize(listWidth, m.height-4)
		m.viewport.Width = m.width - listWidth - 8
		m.viewport.Height = m.height - 22

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			// Esc thì cho thoát luôn cho nhanh nếu đang ở Search hoặc Welcome
			if m.state == StateSearch || m.state == StateWelcome {
				return m, tea.Quit
			}
			// Các màn hình khác thì quay lại
			if m.state == StateEpisodes {
				m.state = StateBrowse
				return m, nil
			}
			if m.state == StateBrowse {
				m.state = StateSearch
				m.textInput.Reset()
				m.textInput.Focus()
				return m, nil
			}

		case "backspace":
			// [QUAN TRỌNG] Nếu đang Search, mình chỉ xoá lỗi (nếu có)
			// Tuyệt đối KHÔNG return m, nil ở đây.
			// Hãy để nó chạy tiếp xuống cuối hàm để m.textInput nhận được phím xoá!
			if m.state == StateSearch {
				if m.err != nil {
					m.err = nil
				}
				// Không return gì cả, để nó "rơi" xuống dưới
			} else {
				// Các màn hình khác thì Backspace vẫn là quay lại
				if m.state == StateWelcome {
					return m, tea.Quit
				}
				if m.state == StateEpisodes {
					m.state = StateBrowse
					return m, nil
				}
				if m.state == StateBrowse {
					m.state = StateSearch
					m.textInput.Reset()
					m.textInput.Focus()
					return m, nil
				}
			}

		case "y", "Y":
			if m.state == StateWelcome {
				return m, playLinkCmd(m.history.LastEpLink)
			}
		case "n", "N":
			if m.state == StateWelcome {
				m.state = StateSearch
				return m, nil
			}

		case "enter":
			if m.state == StateSearch {
				m.err = nil
				return m, searchMoviesCmd(m.textInput.Value())
			} else if m.state == StateBrowse {
				itm, ok := m.movieList.SelectedItem().(Movie)
				if ok {
					return m.enterMovie(itm.Slug)
				}
			} else if m.state == StateEpisodes {
				ep, ok := m.episodeList.SelectedItem().(Episode)
				if ok {
					link := ep.LinkM3u8
					if link == "" {
						link = ep.LinkEmbed
					}
					m.history = History{
						LastMovieName: m.detailCache[m.selectedSlug].Name,
						LastMovieSlug: m.selectedSlug,
						LastEpName:    ep.Name,
						LastEpLink:    link,
						UpdatedAt:     time.Now(),
					}
					saveHistory(m.history)
					return m, playLinkCmd(link)
				}
			}
		}

	case SearchResponse:
		items := []list.Item{}
		for _, v := range msg.Data.Items {
			items = append(items, v)
		}

		// [FIX] Nếu không có kết quả, hiện lỗi chứ không sập!
		if len(items) == 0 {
			m.err = fmt.Errorf("Huhu, không tìm thấy phim này bà ơi! 😢")
			return m, nil
		}

		m.movieList.SetItems(items)
		m.state = StateBrowse
		m.movieList.Select(0)
		cmds = append(cmds, m.triggerDebounceDetail(items[0].(Movie)))

	case DetailLoadedMsg:
		if !msg.IsError {
			m.detailCache[msg.Slug] = msg.Detail
			if m.state == StateBrowse {
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(cleanHTML(msg.Detail.Content)))
			}
			if m.selectedSlug == msg.Slug && m.state == StateBrowse {
				newModel, cmd := m.enterMovie(msg.Slug)
				return newModel, cmd
			}
		}

	case error:
		m.err = msg
	}

	if m.state == StateSearch {
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == StateBrowse {
		curr := m.movieList.SelectedItem()
		m.movieList, cmd = m.movieList.Update(msg)
		cmds = append(cmds, cmd)
		newSel := m.movieList.SelectedItem()
		if curr != nil && newSel != nil && curr.(Movie).Slug != newSel.(Movie).Slug {
			cmds = append(cmds, m.triggerDebounceDetail(newSel.(Movie)))
			m.viewport.SetContent("⏳ Đang tải...")
		}
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.state == StateEpisodes {
		m.episodeList, cmd = m.episodeList.Update(msg)
		cmds = append(cmds, cmd)
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) enterMovie(slug string) (model, tea.Cmd) {
	m.selectedSlug = slug
	detail, ok := m.detailCache[slug]
	if !ok {
		return *m, func() tea.Msg { return fetchDetailMsg(slug, "", true) }
	}

	if len(detail.Episodes) > 0 {
		eps := detail.Episodes[0].ServerData
		items := []list.Item{}
		for _, e := range eps {
			items = append(items, e)
		}
		m.episodeList.SetItems(items)
		m.episodeList.Title = "Chọn tập: " + detail.Name
		m.state = StateEpisodes
	}
	return *m, nil
}

// --- VIEW ---
func (m model) View() string {
	// [FIX] Hiện thông báo lỗi nếu có
	if m.err != nil {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			errorStyle.Render(fmt.Sprintf("🚨 LỖI RỒI MÁ ƠI:\n%v\n\n(Bấm Backspace để quay lại)", m.err)))
	}

	if m.state == StateWelcome {
		content := fmt.Sprintf("Chào cục dzàng quay trở lại! 🥰\n\nLần trước bà đang coi dở:\n%s - Tập %s\n\nBà có mún coi tiếp hem? (Y/N)",
			m.history.LastMovieName, m.history.LastEpName)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, welcomeStyle.Render(content))
	}

	if m.state == StateSearch {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, titleStyle.Render("🎬 RẠP CHIẾU PHIM TẠI GIA"), "\n", m.textInput.View()))
	}

	activeList := m.movieList
	if m.state == StateEpisodes {
		activeList = m.episodeList
	}

	detailView := ""
	curr := m.movieList.SelectedItem()
	if curr != nil {
		slug := curr.(Movie).Slug
		d, found := m.detailCache[slug]
		poster := posterStyle.Render("\n\n  ⏳ Đang tải... \n\n")
		if found {
			poster = posterStyle.Render(d.ASCIIArt)
		}

		detailView = rightPaneStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			movieTitleStyle.Render(curr.(Movie).Name), "\n",
			poster, "\n",
			descStyle.Render(m.viewport.View()),
			"\n"+lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#585b70")).Render("ENTER: Chọn | ESC: Back"),
		))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, appStyle.Render(leftPaneStyle.Render(activeList.View())), detailView)
}

// --- API COMMANDS ---
func searchMoviesCmd(keyword string) tea.Cmd {
	return func() tea.Msg {
		baseURL := os.Getenv("OPHIM_SEARCH_URL")
		if baseURL == "" {
			baseURL = "" //paste vô chỗ này
		}
		res, err := http.Get(baseURL + url.QueryEscape(keyword))
		if err != nil {
			return err
		}
		defer res.Body.Close()
		var r SearchResponse
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return err
		}
		return r
	}
}

func fetchDetailMsg(slug, thumb string, isEnter bool) tea.Msg {
	baseURL := os.Getenv("OPHIM_DETAIL_URL")
	if baseURL == "" {
		baseURL = "" //paste vô chỗ này
	}
	res, err := http.Get(baseURL + slug)
	if err != nil {
		return DetailLoadedMsg{IsError: true}
	}
	defer res.Body.Close()
	var r MovieDetailResponse
	json.NewDecoder(res.Body).Decode(&r)
	d := r.Data.Item

	fallbackBlock := lipgloss.NewStyle().Align(lipgloss.Center).Render(lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f9e2af")).Render("🎬 OPHIM"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Render("(✖╭╮✖)"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#bac2de")).Render("Ảnh lỗi!"),
	))
	d.ASCIIArt = lipgloss.Place(20, 10, lipgloss.Center, lipgloss.Center, fallbackBlock, lipgloss.WithWhitespaceChars(" "))
	return DetailLoadedMsg{Slug: slug, Detail: d}
}

func (m *model) triggerDebounceDetail(movie Movie) tea.Cmd {
	if m.debounceTimer != nil {
		m.debounceTimer.Stop()
	}
	return tea.Tick(DEBOUNCE_DURATION, func(t time.Time) tea.Msg {
		return fetchDetailMsg(movie.Slug, "", false)
	})
}

func playLinkCmd(link string) tea.Cmd {
	c := exec.Command("mpv", link, "--fs")
	return tea.ExecProcess(c, func(err error) tea.Msg { return nil })
}

func cleanHTML(s string) string {
	return strings.NewReplacer("<p>", "", "</p>", "\n", "<br>", "\n", "<strong>", "", "</strong>", "").Replace(s)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
