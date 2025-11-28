package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type post struct {
	Title string
	Desc  string
}

type model struct {
	cb func([]post) tea.Model

	posts []post

	selectionStart int
	selectionSize  int // actually size-1, this way the zero value (0) is ready to use.

	size tea.WindowSizeMsg
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.selectionSize = -1
			fallthrough
		case "space", "enter":
			if m.cb == nil {
				return m, tea.Quit
			}
			rt := m.cb(m.posts[m.selectionStart : m.selectionStart+m.selectionSize+1])
			if rt == nil {
				return m, tea.Quit
			}
			return rt, nil

		case "up", "k":
			m.selectionStart--
		case "down", "j":
			m.selectionStart++
		case "-":
			m.selectionSize--
		case "=":
			m.selectionSize++
		}
	}

	m.selectionStart = max(m.selectionStart, 0)
	m.selectionSize = max(m.selectionSize, 0)

	m.selectionStart = min(m.selectionStart, len(m.posts)-1)
	m.selectionSize = min(m.selectionSize, len(m.posts)-m.selectionStart-1)

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	{
		scrollIndicator := fmt.Sprintf("━━━━ %2d/%d ━━━━", m.selectionStart+1, len(m.posts))
		scrollIndicator = lipgloss.Style{}.
			Foreground(lipgloss.Color("#d0d0d0")).
			Width(m.size.Width).
			AlignHorizontal(lipgloss.Center).
			Render(scrollIndicator)
		b.WriteString(scrollIndicator + "\n")
	}

	const ITEM_HEIGHT = 5

	for offset := range (m.size.Height - 2) / ITEM_HEIGHT {
		actualIndex := offset + m.selectionStart - 1
		var post post
		if 0 <= actualIndex && actualIndex < len(m.posts) {
			post = m.posts[actualIndex]
		}

		selection := "  "
		if m.selectionStart <= actualIndex && actualIndex <= m.selectionStart+m.selectionSize {
			selection = "┃ "
		}

		width := m.size.Width - 4
		item := lipgloss.JoinHorizontal(lipgloss.Center,
			"",
			lipgloss.Style{}.
				Foreground(lipgloss.Color("#808080")).
				Width(width).    // pad to width.
				Height(3).       // pad to height.
				MaxWidth(width). // truncate width if wider.
				MaxHeight(3).    // truncate height if taller.
				Render(post.Desc),
		)
		item = lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.Style{}.
				Bold(true).
				Render(post.Title),
			item,
		)
		item = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.Style{}.
				Foreground(lipgloss.Color("#d0d0d0")).
				Render(strings.Repeat("\n"+selection, ITEM_HEIGHT)[1:]),
			item,
		)
		b.WriteString(item + "\n")
	}

	return b.String()
}

func main() {
	var rt []post
	p := tea.NewProgram(model{
		posts: generatePosts(),
		cb:    func(p []post) tea.Model { rt = p; return nil },
	}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
	for _, p := range rt {
		b, _ := json.MarshalIndent(p, "", "\t")
		b = append(b, '\n')
		os.Stdout.Write(b)
	}
}

func generatePosts() []post {
	return []post{
		{
			Title: "Exploring the Alps",
			Desc:  "A journey through the stunning landscapes of the Swiss Alps, featuring breathtaking views and challenging hikes. This is a very long description that should be truncated.",
		},
		{
			Title: "The Art of Baking",
			Desc:  "Learn the secrets behind delicious pastries and breads from a master baker. This post covers everything from sourdough to croissants. This is a very long description that should be truncated.",
		},
		{
			Title: "A Guide to Urban Gardening",
			Desc:  "Transform your small city space into a green oasis. This guide provides tips on choosing the right plants, containers, and soil for a thriving urban garden. This is a very long description that should be truncated.",
		},
		{
			Title: "The History of Jazz Music",
			Desc:  "Discover the roots of jazz, from its origins in New Orleans to its evolution into a global phenomenon. This post explores the key figures and movements that shaped this iconic genre. This is a very long description that should be truncated.",
		},
		{
			Title: "Mastering the Art of Photography",
			Desc:  "Whether you're a beginner or a seasoned pro, this post offers valuable insights into composition, lighting, and editing to help you take your photography to the next level. This is a very long description that should be truncated.",
		},
		{
			Title: "The Science of Sleep",
			Desc:  "Uncover the mysteries of sleep and learn how to improve your sleep quality. This post delves into the different stages of sleep, the importance of circadian rhythms, and practical tips for a better night's rest. This is a very long description that should be truncated.",
		},
		{
			Title: "A Culinary Tour of Italy",
			Desc:  "Embark on a delicious journey through the diverse culinary regions of Italy. From the fresh seafood of Sicily to the rich pasta dishes of Bologna, this post is a feast for the senses. This is a very long description that should be truncated.",
		},
		{
			Title: "The World of Competitive Gaming",
			Desc:  "Explore the fast-paced and exciting world of eSports. This post takes a look at the most popular games, the top players, and the massive tournaments that draw millions of fans worldwide. This is a very long description that should be truncated.",
		},
		{
			Title: "The Secrets of a Successful Startup",
			Desc:  "Learn from the founders of some of the world's most successful startups. This post shares their stories, strategies, and advice for building a thriving business from the ground up. This is a very long description that should be truncated.",
		},
		{
			Title: "The Ultimate Guide to Fitness",
			Desc:  "Get in the best shape of your life with this comprehensive guide to fitness. This post covers everything from creating a workout plan to proper nutrition and recovery. This is a very long description that should be truncated.",
		},
		{
			Title: "The Beauty of the Night Sky",
			Desc:  "Discover the wonders of the cosmos with this guide to stargazing. Learn how to identify constellations, planets, and other celestial objects with the naked eye or a telescope. This is a very long description that should be truncated.",
		},
		{
			Title: "The Art of Storytelling",
			Desc:  "Master the art of storytelling with this guide to crafting compelling narratives. This post covers the essential elements of a good story, from character development to plot structure. This is a very long description that should be truncated.",
		},
		{
			Title: "The Rise of Artificial Intelligence",
			Desc:  "Explore the fascinating world of artificial intelligence and its impact on our lives. This post examines the latest advancements in AI, from machine learning to natural language processing. This is a very long description that should be truncated.",
		},
		{
			Title: "The Magic of a Good Book",
			Desc:  "Rediscover the joy of reading with this celebration of the written word. This post explores the power of books to transport us to new worlds, challenge our perspectives, and enrich our lives. This is a very long description that should be truncated.",
		},
		{
			Title: "The Thrill of Adventure Travel",
			Desc:  "Embark on an unforgettable adventure with this guide to thrill-seeking travel. From bungee jumping in New Zealand to white-water rafting in Costa Rica, this post is your ticket to an adrenaline-fueled journey. This is a very long description that should be truncated.",
		},
		{
			Title: "The Power of Positive Thinking",
			Desc:  "Transform your life with the power of positive thinking. This post explores the science behind optimism and provides practical strategies for cultivating a more positive mindset. This is a very long description that should be truncated.",
		},
		{
			Title: "The History of the Internet",
			Desc:  "Journey back in time to the early days of the internet and discover how it evolved into the global network we know today. This post covers the key milestones, technologies, and people that shaped the digital age. This is a very long description that should be truncated.",
		},
		{
			Title: "The Joy of Cooking for Others",
			Desc:  "Experience the satisfaction of sharing a home-cooked meal with loved ones. This post offers tips and recipes for hosting a memorable dinner party, from planning the menu to creating a warm and inviting atmosphere. This is a very long description that should be truncated.",
		},
		{
			Title: "The Future of Space Exploration",
			Desc:  "Look to the stars and imagine the future of space exploration. This post examines the latest missions, technologies, and discoveries that are pushing the boundaries of our understanding of the universe. This is a very long description that should be truncated.",
		},
		{
			Title: "The Importance of Lifelong Learning",
			Desc:  "Embrace the pursuit of knowledge with this celebration of lifelong learning. This post explores the benefits of continuous learning, from personal growth to professional development, and provides tips for staying curious and engaged throughout your life. This is a very long description that should be truncated.",
		},
	}
}
