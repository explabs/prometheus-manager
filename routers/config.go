package routers

import (
	"fmt"
	"net/http"
)

type PrometheusConfig struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		ScrapeTimeout      string `yaml:"scrape_timeout"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	ScrapeConfigs []ScrapeConfigs `yaml:"scrape_configs"`
}
type ScrapeConfigs struct {
	JobName       string `yaml:"job_name"`
	StaticConfigs []struct {
		Targets []string `yaml:"targets"`
	} `yaml:"static_configs"`
	MetricsPath string `yaml:"metrics_path,omitempty"`
	BasicAuth   struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"basic_auth,omitempty"`
	ScrapeInterval string `yaml:"scrape_interval,omitempty"`
	ScrapeTimeout  string `yaml:"scrape_timeout,omitempty"`
}
type Jobs struct {
	Jobs     []JsonData `json:"jobs"`
	Password string     `json:"password"`
	Target   string     `json:"target"`
	Interval string     `json:"interval"`
	Timeout  string     `json:"timeout"`
}
type JsonData struct {
	Name 	 string	`json:"name"`
	Path     string `json:"path"`
	Interval string `json:"interval"`
}

func (j *Jobs) ConvertToYml() error {
	var p *PrometheusConfig
	p.Global.ScrapeInterval = j.Interval
	p.Global.ScrapeTimeout = j.Timeout
	p.Global.EvaluationInterval = j.Interval

	for _, job := range j.Jobs {
		scrapeConfig := ScrapeConfigs{
			JobName:     job.Name,
			StaticConfigs: []struct {
				Targets []string `yaml:"targets"`
			}{},
			MetricsPath: job.Path,
			BasicAuth: struct {
				Username string `yaml:"username"`
				Password string `yaml:"password"`
			}{
				Username: "checker",
				Password: j.Password,
			},
			ScrapeInterval: job.Interval,
			ScrapeTimeout:  j.Timeout,
		}
		fmt.Println(scrapeConfig)
	}

	p.GenerateConfig("")
	return nil
}

func (p *PrometheusConfig) GenerateConfig(filepath string) error {
	fmt.Println(filepath, p)
	return nil
}

func JsonParse(w http.ResponseWriter, r *http.Request){
	
}