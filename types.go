//    Copyright 2017 Yoshi Yamaguchi
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

// https://dialogflow.com/docs/fulfillment#request
type Request struct {
	Lang            string           `json:"lang"`
	Status          *RequestStatus   `json:"status"`
	Timestamp       string           `json:"timestamp"`
	SessionID       string           `json:"sessionId"`
	Result          *RequestResult   `json:"result"`
	ID              string           `json:"id"`
	OriginalRequest *OriginalRequest `json:"originalRequest"`
}

func (req *Request) GetIntent() string {
	return req.Result.Action
}

type RequestStatus struct {
	ErrorType string `json:"errorType"`
	Code      int    `json:"code"`
}

type RequestResult struct {
	Parameters       map[string]interface{} `json:"parameters"`
	Contexts         []*Context             `json:"contexts"`
	ResolvedQuery    string                 `json:"resolvedQuery"`
	Source           string                 `json:"source"`
	Score            float32                `json:"score"`
	Speech           string                 `json:"speech"`
	Fulfillment      *RequestFulfillment    `json:"fulfillment"`
	ActionIncomplete bool                   `json:"actionIncomplete"`
	Action           string                 `json:"action"`
	Metadata         *Metadata              `json:"metadata"`
}

type Context struct {
	ContextName string                 `json:"name"`
	LifeSpan    int                    `json:"lifespan"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type RequestFulfillment struct {
	Messages []*Message `json:"messages"`
	Speech   string     `json:"speech"`
}

type Message struct {
	Speech string `json:"speech"`
	Type   int    `json:"type"`
}

// TODO: check if there are other values than "true" and "false" for webhookUsed and webhookForSlogFillingUsed.
type Metadata struct {
	IntentID                  string `json:"intentId"`
	WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
	IntentName                string `json:"intentName"`
	WebhookUsed               string `json:"webhookUsed"`
}

type OriginalRequest struct {
	Source string               `json:"source"`
	Data   *OriginalRequestData `json:"data"`
}

type OriginalRequestData struct {
	Inputs       []*InputData  `json:"inputs"`
	User         *User         `json:"user"`
	Conversation *Conversation `json:"conversation"`
}

type InputData struct {
	RawInputs []*RawInput  `json:"raw_inputs"`
	Intent    string       `json:"intent"`
	Arguments []*Arguments `json:"arguments"`
}

type RawInput struct {
	Query     string `json:"query"`
	InputType int    `json:"input_type"`
}

type Arguments struct {
	TextValue string `json:"text_value"`
	RawText   string `json:"raw_text"`
	Name      string `json:"name"`
}

type User struct {
	UserID string `json:"user_id"`
}

type Conversation struct {
	ConversationID    string      `json:"conversation_id"`
	Type              interface{} `json:"type"`
	ConversationToken string      `json:"conversation_token"`
}

// https://dialogflow.com/docs/fulfillment#response
// TODO: find the good method to add `data` field.
type Response struct {
	Speech        string         `json:"speech"`
	DisplayText   string         `json:"displayText"`
	ContextOut    []*Context     `json:"contextOut"`
	Source        string         `json:"source"`
	FollowupEvent *FollowupEvent `json:"followupEvent"`
}

func NewResponse(speech string) *Response {
	return &Response{
		Speech: speech,
	}
}

type FollowupEvent struct {
	EventName string                 `json:"name"`
	Data      map[string]interface{} `json:"data"`
}

func (res *Response) SetDisplayText(text string) *Response {
	res.DisplayText = text
	return res
}
