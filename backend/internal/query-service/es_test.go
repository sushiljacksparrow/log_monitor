package queryservice

// func timePtr(t time.Time) *time.Time { return &t }
// func TestSearchAuthLogs(t *testing.T) {
// 	config, err := config.InitConfig()
// 	if err != nil {
// 		log.Fatalf("error while init envs: %v", config)
// 	}
// 	esClient, esTypedClient, err := elasticsearch.InitES(config)
// 	if err != nil {
// 		log.Fatalf("error while init es client", err)
// 	}
// 	if err != nil {
// 		log.Fatalf("error while init bulkindexer: %v", err)
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	body := map[string]string{
// 		"service":    "auth-service",
// 		"level":      "WARN",
// 		"message":    "Invalid password attempt",
// 		"request_id": "550e8400-e29b-41d4-a716-446655440000",
// 		"user_id":    "980e1209-e29b-41d4-a716-446655440000",
// 		"ip":         "192.168.1.10",
// 		"timestamp":  time.Now().Format(time.RFC3339),
// 	}
// 	// 🔥 Convert map → JSON
// 	jsonBody, err := json.Marshal(body)
// 	if err != nil {
// 		log.Fatalf("error while converting into json %v", err)
// 	}

// 	res, err := esClient.Index(constants.AUTH_SERVICE_LOGS_INDEX, bytes.NewReader(jsonBody))
// 	if err != nil {
// 		log.Fatalf("error while indexing: %v", err)
// 	}
// 	defer res.Body.Close()
// 	repo := NewRepository(esClient, esTypedClient)

// 	now := time.Now()

// 	tests := []struct {
// 		name     string
// 		authlog  AuthLogFilters
// 		noOfDocs int
// 	}{
// 		{
// 			name: "only request_id valid",
// 			authlog: AuthLogFilters{
// 				RequestID: strPtr("550e8400-e29b-41d4-a716-446655440000"),
// 			},
// 			noOfDocs: 1,
// 		},

// 		{
// 			name: "level and message valid",
// 			authlog: AuthLogFilters{
// 				Level:   strPtr("WARN"),
// 				Message: strPtr("Invalid password attempt"),
// 			},
// 			noOfDocs: ,
// 		},

// 		{
// 			name: "user_id and ip valid",
// 			authlog: AuthLogFilters{
// 				UserID: strPtr("980e1209-e29b-41d4-a716-446655440000"),
// 				IP:     strPtr("192.168.1.10"),
// 			},
// 		},

// 		{
// 			name: "all fields valid",
// 			authlog: AuthLogFilters{
// 				RequestID:      strPtr("550e8400-e29b-41d4-a716-446655440000"),
// 				Level:          strPtr("WARN"),
// 				Message:        strPtr("Invalid password attempt"),
// 				UserID:         strPtr("980e1209-e29b-41d4-a716-446655440000"),
// 				IP:             strPtr("192.168.1.10"),
// 				StartTimestamp: timestamppb.New(now),
// 				EndTimestamp:   timestamppb.New(now),
// 			},
// 		},

// 		{
// 			name:    "empty filter",
// 			authlog: AuthLogFilters{},
// 		},

// 		// ❌ WRONG / EDGE CASES

// 		{
// 			name: "wrong request_id",
// 			authlog: AuthLogFilters{
// 				RequestID: strPtr("wrong-id"),
// 			},
// 		},

// 		{
// 			name: "wrong level",
// 			authlog: AuthLogFilters{
// 				Level: strPtr("DEBUG"),
// 			},
// 		},

// 		{
// 			name: "wrong message",
// 			authlog: AuthLogFilters{
// 				Message: strPtr("Some random message"),
// 			},
// 		},

// 		{
// 			name: "wrong user_id",
// 			authlog: AuthLogFilters{
// 				UserID: strPtr("11111111-e29b-41d4-a716-446655440000"),
// 			},
// 		},

// 		{
// 			name: "wrong ip",
// 			authlog: AuthLogFilters{
// 				IP: strPtr("10.0.0.1"),
// 			},
// 		},

// 		{
// 			name: "future timestamp (no match)",
// 			authlog: AuthLogFilters{
// 				StartTimestamp: timestamppb.New(now.Add(24 * time.Hour)),
// 			},
// 		},

// 		{
// 			name: "mixed valid and invalid",
// 			authlog: AuthLogFilters{
// 				RequestID: strPtr("550e8400-e29b-41d4-a716-446655440000"),
// 				IP:        strPtr("10.0.0.1"),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			res, err := repo.SearchAuthLogs(ctx, tt.authlog)
// 			if err != nil {
// 				t.Fatalf("error while executing test case: %s - %v,", tt.name, err)
// 			}
// 			for _, doc := range res {
// 				doc
// 			}
// 		})
// 	}
// }
