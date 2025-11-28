package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type post struct {
	title       string
	description string
}

type model struct {
	cb func([]post) tea.Model

	posts []post

	selectionStart int
	selectionSize  int // actually size-1, this way the zero value (0) is ready to use.

	height         int
	viewportOffset int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
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
			if m.selectionStart > 0 {
				m.selectionStart--
				if m.selectionStart < m.viewportOffset {
					m.viewportOffset = m.selectionStart
				}
			}
		case "down", "j":
			if m.selectionStart < len(m.posts)-1 {
				m.selectionStart++
				if m.selectionStart >= m.viewportOffset+m.height/4 {
					m.viewportOffset = m.selectionStart - m.height/4 + 1
				}
			}
		case "-":
			if m.selectionSize > 0 {
				m.selectionSize--
			}
		case "=":
			m.selectionSize++
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Scroll indicator
	b.WriteString(fmt.Sprintf("\n---\n%d/%d", m.selectionStart+1, len(m.posts)))

	// Determine the slice of posts to render
	start := m.viewportOffset
	end := min(m.viewportOffset+m.height/4, len(m.posts))
	visiblePosts := m.posts[start:end]

	for i, post := range visiblePosts {
		actualIndex := start + i
		if m.selectionStart <= actualIndex && actualIndex <= m.selectionStart+m.selectionSize {
			b.WriteString("* ")
		} else {
			b.WriteString("  ")
		}

		b.WriteString(post.title)
		b.WriteString("\n")
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Center,
			"    ",
			truncate(post.description, 3),
		))
		b.WriteString("\n\n")
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
		fmt.Println(p.title)
	}
}

func truncate(s string, maxLines int) string {
	lines := strings.Split(s, "\n")
	if len(lines) > maxLines {
		return strings.Join(lines[:maxLines], "\n") + "..."
	}
	return s
}

func generatePosts() []post {
	return []post{
		{
			title:       "Exploring the Alps",
			description: "A journey through the stunning landscapes of the Swiss Alps, featuring breathtaking views and challenging hikes.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Art of Baking",
			description: "Learn the secrets behind delicious pastries and breads from a master baker.\nThis post covers everything from sourdough to croissants.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "A Guide to Urban Gardening",
			description: "Transform your small city space into a green oasis.\nThis guide provides tips on choosing the right plants, containers, and soil for a thriving urban garden.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The History of Jazz Music",
			description: "Discover the roots of jazz, from its origins in New Orleans to its evolution into a global phenomenon.\nThis post explores the key figures and movements that shaped this iconic genre.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "Mastering the Art of Photography",
			description: "Whether you're a beginner or a seasoned pro, this post offers valuable insights into composition, lighting, and editing to help you take your photography to the next level.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Science of Sleep",
			description: "Uncover the mysteries of sleep and learn how to improve your sleep quality.\nThis post delves into the different stages of sleep, the importance of circadian rhythms, and practical tips for a better night's rest.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "A Culinary Tour of Italy",
			description: "Embark on a delicious journey through the diverse culinary regions of Italy.\nFrom the fresh seafood of Sicily to the rich pasta dishes of Bologna, this post is a feast for the senses.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The World of Competitive Gaming",
			description: "Explore the fast-paced and exciting world of eSports.\nThis post takes a look at the most popular games, the top players, and the massive tournaments that draw millions of fans worldwide.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Secrets of a Successful Startup",
			description: "Learn from the founders of some of the world's most successful startups.\nThis post shares their stories, strategies, and advice for building a thriving business from the ground up.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Ultimate Guide to Fitness",
			description: "Get in the best shape of your life with this comprehensive guide to fitness.\nThis post covers everything from creating a workout plan to proper nutrition and recovery.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Beauty of the Night Sky",
			description: "Discover the wonders of the cosmos with this guide to stargazing.\nLearn how to identify constellations, planets, and other celestial objects with the naked eye or a telescope.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Art of Storytelling",
			description: "Master the art of storytelling with this guide to crafting compelling narratives.\nThis post covers the essential elements of a good story, from character development to plot structure.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Rise of Artificial Intelligence",
			description: "Explore the fascinating world of artificial intelligence and its impact on our lives.\nThis post examines the latest advancements in AI, from machine learning to natural language processing.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Magic of a Good Book",
			description: "Rediscover the joy of reading with this celebration of the written word.\nThis post explores the power of books to transport us to new worlds, challenge our perspectives, and enrich our lives.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Thrill of Adventure Travel",
			description: "Embark on an unforgettable adventure with this guide to thrill-seeking travel.\nFrom bungee jumping in New Zealand to white-water rafting in Costa Rica, this post is your ticket to an adrenaline-fueled journey.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Power of Positive Thinking",
			description: "Transform your life with the power of positive thinking.\nThis post explores the science behind optimism and provides practical strategies for cultivating a more positive mindset.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The History of the Internet",
			description: "Journey back in time to the early days of the internet and discover how it evolved into the global network we know today.\nThis post covers the key milestones, technologies, and people that shaped the digital age.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Joy of Cooking for Others",
			description: "Experience the satisfaction of sharing a home-cooked meal with loved ones.\nThis post offers tips and recipes for hosting a memorable dinner party, from planning the menu to creating a warm and inviting atmosphere.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Future of Space Exploration",
			description: "Look to the stars and imagine the future of space exploration.\nThis post examines the latest missions, technologies, and discoveries that are pushing the boundaries of our understanding of the universe.\nThis is a very long description that should be truncated.",
		},
		{
			title:       "The Importance of Lifelong Learning",
			description: "Embrace the pursuit of knowledge with this celebration of lifelong learning.\nThis post explores the benefits of continuous learning, from personal growth to professional development, and provides tips for staying curious and engaged throughout your life.\nThis is a very long description that should be truncated.",
		},
	}
}
