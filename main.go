package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"

    "github.com/firebase/genkit/go/genkit"
    "github.com/firebase/genkit/go/plugins/googlegenai"
    "github.com/firebase/genkit/go/plugins/server"
    "github.com/joho/godotenv"
)

// Define input schema
type RecipeInput struct {
    Ingredient           string `json:"ingredient" jsonschema:"description=Main ingredient or cuisine type"`
    DietaryRestrictions  string `json:"dietaryRestrictions,omitempty" jsonschema:"description=Any dietary restrictions"`
}

// Define output schema
type Recipe struct {
    Title        string   `json:"title"`
    Description  string   `json:"description"`
    PrepTime     string   `json:"prepTime"`
    CookTime     string   `json:"cookTime"`
    Servings     int      `json:"servings"`
    Ingredients  []string `json:"ingredients"`
    Instructions []string `json:"instructions"`
    Tips         []string `json:"tips,omitempty"`
}

// wrapBodyForGenkit ensures the request body has a top-level "data" key so Genkit's
// handler receives the flow input. Accepts both {"data": {...}} and raw {...}.
func wrapBodyForGenkit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || r.ContentLength == 0 {
			next.ServeHTTP(w, r)
			return
		}
		raw, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var parsed map[string]json.RawMessage
		if err := json.Unmarshal(raw, &parsed); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, hasData := parsed["data"]; hasData {
			r.Body = io.NopCloser(bytes.NewReader(raw))
		} else {
			wrapped, _ := json.Marshal(map[string]json.RawMessage{"data": raw})
			r.Body = io.NopCloser(bytes.NewReader(wrapped))
			r.ContentLength = int64(len(wrapped))
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
    // Load .env so GEMINI_API_KEY is set when run via genkit start
    _ = godotenv.Load()

    ctx := context.Background()

    // Initialize Genkit with the Google AI plugin and load prompts from ./prompts
    g := genkit.Init(ctx,
        genkit.WithPlugins(&googlegenai.GoogleAI{}),
        genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
        genkit.WithPromptDir("./prompts"),
    )

    // Define a recipe generator flow
    recipeGeneratorFlow := genkit.DefineFlow(g, "recipeGeneratorFlow", func(ctx context.Context, input *RecipeInput) (*Recipe, error) {
        // Create a prompt based on the input
        dietaryRestrictions := input.DietaryRestrictions
        if dietaryRestrictions == "" {
            dietaryRestrictions = "none"
        }





		        // prompt := fmt.Sprintf(`Create a recipe with the following requirements:
        //     Main ingredient: %s
        //     Dietary restrictions: %s`, input.Ingredient, dietaryRestrictions)

		// Look up a .prompt file with type information
		recipePrompt := genkit.LookupDataPrompt[RecipeInput, *Recipe](g, "create_recipe")

		// Execute with strongly-typed input, get strongly-typed output
		recipe, _, err := recipePrompt.Execute(ctx, RecipeInput{Ingredient: input.Ingredient, DietaryRestrictions: dietaryRestrictions})
		if err != nil {
			return nil, err
		}

		return recipe, nil
    })

    // Run the flow once to test it
    recipe, err := recipeGeneratorFlow.Run(ctx, &RecipeInput{
        Ingredient:          "avocado",
        DietaryRestrictions: "vegetarian",
    })
    if err != nil {
        log.Fatalf("could not generate recipe: %v", err)
    }

    // Print the structured recipe
    recipeJSON, _ := json.MarshalIndent(recipe, "", "  ")
    fmt.Println("Sample recipe generated:")
    fmt.Println(string(recipeJSON))

    // Start a server to serve the flow and keep the app running for the Developer UI
    mux := http.NewServeMux()
    flowHandler := genkit.Handler(recipeGeneratorFlow)
    mux.HandleFunc("POST /recipeGeneratorFlow", wrapBodyForGenkit(flowHandler))

    log.Println("Starting server on http://localhost:3400")
    log.Println("Flow available at: POST http://localhost:3400/recipeGeneratorFlow")
    log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}