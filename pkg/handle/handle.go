package handle

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atompi/mpsbot/pkg/options"
	redisutil "github.com/atompi/mpsbot/pkg/util/redis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type TargetYaml struct {
	Targets []string       `yaml:"targets"`
	Labels  map[string]int `yaml:"labels"`
}

func Handle(opts options.Options) {
	c := redisutil.New(opts.Redis)

	for {
		iter := scan(c, opts.Redis.Prefix, opts.Redis.DialTimeout)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.Redis.DialTimeout)*time.Second)
		data := make(map[string][]string)
		for iter.Next(ctx) {
			key := iter.Val()

			s := strings.TrimPrefix(key, opts.Redis.Prefix)
			parts := strings.Split(s, "__")
			if len(parts) != 2 {
				zap.S().Warnf("invalid key: %s", key)
				continue
			}

			module := parts[0]
			instance := parts[1]
			data[module] = append(data[module], instance)
		}

		if err := iter.Err(); err != nil {
			zap.S().Errorf("scan failed: %v", err)
			continue
		}
		cancel()

		err := writeToYAML(opts.Task.OutputPath, data, opts.Task.VmTenantLabels)
		if err != nil {
			zap.S().Errorf("write to yaml failed: %v", err)
		}

		time.Sleep(time.Duration(opts.Task.Interval) * time.Second)
	}
}

func scan(c *redis.Client, prefix string, timeout int) *redis.ScanIterator {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	iter := c.Scan(ctx, 0, prefix+"*", 0).Iterator()
	cancel()
	return iter
}

func writeToYAML(path string, data map[string][]string, labels map[string]int) error {
	for k, v := range data {
		fileName := fmt.Sprintf("%s.yaml", k)
		filePath := filepath.Join(path, fileName)
		raw := []TargetYaml{}
		t := TargetYaml{
			Targets: v,
			Labels:  labels,
		}
		raw = append(raw, t)

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer file.Close()

		yamlData, err := yaml.Marshal(raw)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %v", err)
		}

		_, err = file.Write(yamlData)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}
