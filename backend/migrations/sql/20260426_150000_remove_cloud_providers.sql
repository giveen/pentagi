-- +goose Up
-- +goose StatementBegin
-- Remove cloud providers: anthropic, gemini, bedrock, deepseek, glm, kimi, qwen
-- Keep only: openai, ollama, custom
DELETE FROM providers WHERE type IN ('anthropic', 'gemini', 'bedrock', 'deepseek', 'glm', 'kimi', 'qwen');
DELETE FROM flows WHERE model_provider_type IN ('anthropic', 'gemini', 'bedrock', 'deepseek', 'glm', 'kimi', 'qwen');
DELETE FROM assistants WHERE model_provider_type IN ('anthropic', 'gemini', 'bedrock', 'deepseek', 'glm', 'kimi', 'qwen');

-- Create new enum type with only 3 providers
CREATE TYPE PROVIDER_TYPE_NEW AS ENUM (
  'openai',
  'ollama',
  'custom'
);

-- Update columns to use the new enum type
ALTER TABLE providers
    ALTER COLUMN type TYPE PROVIDER_TYPE_NEW USING type::text::PROVIDER_TYPE_NEW;

ALTER TABLE flows
    ALTER COLUMN model_provider_type TYPE PROVIDER_TYPE_NEW USING model_provider_type::text::PROVIDER_TYPE_NEW;

ALTER TABLE assistants
    ALTER COLUMN model_provider_type TYPE PROVIDER_TYPE_NEW USING model_provider_type::text::PROVIDER_TYPE_NEW;

-- Drop the old type and rename the new one
DROP TYPE PROVIDER_TYPE;
ALTER TYPE PROVIDER_TYPE_NEW RENAME TO PROVIDER_TYPE;

-- Ensure NOT NULL constraints are preserved
ALTER TABLE providers
    ALTER COLUMN type SET NOT NULL;

ALTER TABLE flows
    ALTER COLUMN model_provider_type SET NOT NULL;

ALTER TABLE assistants
    ALTER COLUMN model_provider_type SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert to enum with all 10 providers
CREATE TYPE PROVIDER_TYPE_NEW AS ENUM (
  'openai',
  'anthropic',
  'gemini',
  'bedrock',
  'ollama',
  'custom',
  'deepseek',
  'glm',
  'kimi',
  'qwen'
);

-- Update columns to use the reverted enum type
ALTER TABLE providers
    ALTER COLUMN type TYPE PROVIDER_TYPE_NEW USING type::text::PROVIDER_TYPE_NEW;

ALTER TABLE flows
    ALTER COLUMN model_provider_type TYPE PROVIDER_TYPE_NEW USING model_provider_type::text::PROVIDER_TYPE_NEW;

ALTER TABLE assistants
    ALTER COLUMN model_provider_type TYPE PROVIDER_TYPE_NEW USING model_provider_type::text::PROVIDER_TYPE_NEW;

-- Drop the new type and rename the old one
DROP TYPE PROVIDER_TYPE;
ALTER TYPE PROVIDER_TYPE_NEW RENAME TO PROVIDER_TYPE;

-- Ensure NOT NULL constraints are preserved
ALTER TABLE providers
    ALTER COLUMN type SET NOT NULL;

ALTER TABLE flows
    ALTER COLUMN model_provider_type SET NOT NULL;

ALTER TABLE assistants
    ALTER COLUMN model_provider_type SET NOT NULL;
-- +goose StatementEnd
