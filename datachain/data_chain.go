package datachain

import (
	"context"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools/sqldatabase"
)

const _DECIDER_TEMPLATE = `Given the below input question and list of potential tables, output a comma separated list of table names that may be necessary to answer this question.

Question: {{.query}}

Table Names: {{.table_names}}

Relevant Table Names:`

type DataChain struct {
	deciderChain *chains.LLMChain
	outputParser schema.OutputParser[[]string]
	SQLChain     *chains.SQLDatabaseChain
}

// New returns a new TableChain.
// useAllTables: user all tables or choose by query
// If you have a lot of tables, you'd better filter out the tables by useAllTables=false,
// otherwise, the token may execeded token limit .
func New(model llms.LanguageModel, engine, dsn string, useAllTables bool) (*DataChain, error) {
	var deciderChain *chains.LLMChain
	if !useAllTables {
		deciderChain = chains.NewLLMChain(model,
			prompts.NewPromptTemplate(_DECIDER_TEMPLATE, []string{"query", "table_names"}))
	}
	db, err := sqldatabase.NewSQLDatabaseWithDSN(engine, dsn, nil)
	if err != nil {
		return nil, err
	}
	sqlChain := chains.NewSQLDatabaseChain(model, 10, db)

	return &DataChain{
		deciderChain: deciderChain,
		SQLChain:     sqlChain,
		outputParser: outputparser.NewCommaSeparatedList(),
	}, nil
}

func (c *DataChain) Call(ctx context.Context, values map[string]any, options ...chains.ChainCallOption) (map[string]any, error) {
	if c.deciderChain != nil {
		values["table_names"] = c.SQLChain.Database.TableNames()
		text, err := chains.Predict(ctx, c.deciderChain, values, options...)
		if err != nil {
			return nil, err
		}
		tbs, err := c.outputParser.Parse(text)
		if err != nil {
			return nil, err
		}
		if len(tbs) > 0 {
			values["table_names_to_use"] = tbs
		}
	}

	ret, err := chains.Predict(ctx, c.SQLChain, values, options...)
	if err != nil {
		return nil, err
	}
	values[c.SQLChain.OutputKey] = ret
	return values, nil
}

// GetMemory returns the memory.
func (c *DataChain) GetMemory() schema.Memory { //nolint:ireturn
	return memory.NewSimple() //nolint:ireturn
}

// GetInputKeys returns the expected input keys.
func (c *DataChain) GetInputKeys() []string {
	return append([]string{}, "query")
}

// GetOutputKeys returns the output keys the chain will return.
func (c *DataChain) GetOutputKeys() []string {
	return c.SQLChain.GetOutputKeys()
}

func (c *DataChain) Close() error {
	return c.SQLChain.Database.Close()
}

func (c *DataChain) Run(ctx context.Context, query string) (string, []string, error) {
	out, err := chains.Call(ctx, c, map[string]any{c.GetInputKeys()[0]: query})
	if err != nil {
		return "", nil, err
	}
	return out[c.GetOutputKeys()[0]].(string), out["table_names_to_use"].([]string), nil
}
