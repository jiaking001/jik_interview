package ai

import (
	"fmt"
	"testing"
)

func TestAI_doChat(t *testing.T) {
	type fields struct {
		systemPrompt string
		userPrompt   string
		model        string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				systemPrompt: "你是一个智能助手",
				userPrompt:   "你好介绍一下你自己",
			},
		},
		{
			name: "test2",
			fields: fields{
				systemPrompt: "你是一个智能助手",
				userPrompt:   "你好",
				model:        "deepseek-v3-250324",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(DoChat(tt.fields.systemPrompt, tt.fields.userPrompt, tt.fields.model))
			//if got := a.doChat(); got != tt.want {
			//	t.Errorf("doChat() = %v, want %v", got, tt.want)
			//}
		})
	}
}
