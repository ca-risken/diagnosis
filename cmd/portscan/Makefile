.PHONY: all run
all: run

run:
	ls *.go | grep -v '_test.go' | xargs go run

enque:
	aws sqs send-message \
		--endpoint-url http://localhost:9324 \
		--queue-url http://localhost:9324/queue/diagnosis-portscan \
		--message-body '{"data_source":"diagnosis:portscan", "project_id":1001, "portscan_setting_id":1001, "portscan_target_id":1001, "target":"127.0.0.1", "scan_only":"false"}'
