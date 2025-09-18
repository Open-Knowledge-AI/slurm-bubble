package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Job struct {
	ID        string
	User      string
	Partition string
	State     string
	Time      string
	Nodes     string
	Name      string
}

type GpuUsage struct {
	Node        string
	GpuIndex    int
	MemUsedMB   int
	MemTotalMB  int
	Utilization int
	JobID       string
}

type viewMode int

const (
	jobList viewMode = iota
	jobDetail
	gpuCluster
	help
)

type model struct {
	jobs        []Job
	gpus        []GpuUsage
	selectedIdx int
	mode        viewMode
	loading     bool
	errorMsg    string
	table       table.Model
}

func initialModel() model {
	columns := []table.Column{
		{Title: "ID", Width: 6},
		{Title: "User", Width: 10},
		{Title: "Partition", Width: 10},
		{Title: "State", Width: 8},
		{Title: "Time", Width: 8},
		{Title: "Nodes", Width: 6},
		{Title: "Name", Width: 20},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
	)
	return model{
		jobs:        []Job{},
		gpus:        []GpuUsage{},
		selectedIdx: 0,
		mode:        jobList,
		table:       t,
		loading:     true,
		errorMsg:    "",
	}
}

// Simulate async loading of jobs
func loadJobsCmd() bubbletea.Cmd {
	return func() bubbletea.Msg {
		time.Sleep(1 * time.Second) // simulate delay
		jobs := []Job{
			{"123", "alice", "gpu", "R", "00:10:00", "2", "train-model"},
			{"124", "bob", "cpu", "PD", "00:00:00", "1", "data-prep"},
		}
		return jobsLoadedMsg(jobs)
	}
}

type jobsLoadedMsg []Job
type gpuLoadedMsg []GpuUsage

func (m model) Init() bubbletea.Cmd {
	return loadJobsCmd()
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {

	case jobsLoadedMsg:
		m.jobs = msg
		m.loading = false
		rows := make([]table.Row, len(m.jobs))
		for i, j := range m.jobs {
			rows[i] = table.Row{j.ID, j.User, j.Partition, j.State, j.Time, j.Nodes, j.Name}
		}
		m.table.SetRows(rows)
		return m, nil

	case bubbletea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, bubbletea.Quit
		case "r":
			m.loading = true
			return m, loadJobsCmd()
		case "g":
			if m.mode == gpuCluster {
				m.mode = jobList
			} else {
				m.mode = gpuCluster
				// Placeholder: could trigger GPU load here
			}
			return m, nil
		case "up":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down":
			if m.selectedIdx < len(m.jobs)-1 {
				m.selectedIdx++
			}
		case "enter":
			m.mode = jobDetail
			return m, nil
		case "esc":
			m.mode = jobList
			return m, nil
		}
	}
	return m, nil
}

func renderJobList(m model) string {
	if m.loading {
		return "Loading jobs..."
	}
	return m.table.View()
}

func renderJobDetail(m model) string {
	if len(m.jobs) == 0 {
		return "No job selected."
	}
	job := m.jobs[m.selectedIdx]
	return fmt.Sprintf(
		"Job Detail\nID: %s\nUser: %s\nPartition: %s\nState: %s\nTime: %s\nNodes: %s\nName: %s\n\nPress ESC to go back",
		job.ID, job.User, job.Partition, job.State, job.Time, job.Nodes, job.Name,
	)
}

func renderGpuView(m model) string {
	if len(m.gpus) == 0 {
		return "No GPU data available."
	}
	out := lipgloss.NewStyle().Bold(true).Render("Cluster GPU Usage\n")
	for _, g := range m.gpus {
		bar := renderProgressBar(g.MemUsedMB, g.MemTotalMB, 20)
		out += fmt.Sprintf("Node: %s GPU%d [%s] %d%% Job: %s\n", g.Node, g.GpuIndex, bar, g.Utilization, g.JobID)
	}
	return out
}

func renderProgressBar(used, total, width int) string {
	fill := int(float64(used) / float64(total) * float64(width))
	bar := ""
	for i := 0; i < fill; i++ {
		bar += "█"
	}
	for i := fill; i < width; i++ {
		bar += "─"
	}
	return bar
}

func (m model) View() string {
	switch m.mode {
	case jobList:
		return renderJobList(m) + "\nPress 'g' for GPU view, 'r' to refresh, 'q' to quit."
	case jobDetail:
		return renderJobDetail(m)
	case gpuCluster:
		return renderGpuView(m) + "\nPress 'g' to go back to jobs."
	default:
		return "Unknown mode"
	}
}

func main() {
	m := initialModel()
	p := bubbletea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
