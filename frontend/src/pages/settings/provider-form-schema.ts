import { z } from 'zod';

// Define agent configuration schema
const agentConfigSchema = z
    .object({
        frequencyPenalty: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        maxLength: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        maxTokens: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        minLength: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        model: z.preprocess((value) => value || '', z.string().min(1, 'Model is required')),
        presencePenalty: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        price: z
            .object({
                cacheRead: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.number().nullable().optional(),
                ),
                cacheWrite: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.number().nullable().optional(),
                ),
                input: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.number().nullable().optional(),
                ),
                output: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.number().nullable().optional(),
                ),
            })
            .nullable()
            .optional(),
        reasoning: z
            .object({
                effort: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.string().nullable().optional(),
                ),
                maxTokens: z.preprocess(
                    (value) => (value === '' || value === undefined ? null : value),
                    z.number().nullable().optional(),
                ),
            })
            .nullable()
            .optional(),
        repetitionPenalty: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        temperature: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        topK: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
        topP: z.preprocess(
            (value) => (value === '' || value === undefined ? null : value),
            z.number().nullable().optional(),
        ),
    })
    .optional();

// Define form schema
const formSchema = z.object({
    agents: z.record(z.string(), agentConfigSchema).optional(),
    apiKey: z.preprocess((value) => value || '', z.string().optional()),
    apiUrl: z.preprocess(
        (value) => (value === '' || value === undefined || value === null ? undefined : value),
        z.string().url('Must be a valid URL').optional(),
    ),
    name: z.preprocess(
        (value) => value || '',
        z.string().min(1, 'Provider name is required').max(50, 'Maximum 50 characters allowed'),
    ),
    type: z.preprocess((value) => value || '', z.string().min(1, 'Provider type is required')),
});

export type EndpointHealth = {
    error?: string;
    model?: string;
    models?: number;
    reachable: boolean;
    statusCode?: number;
    url: string;
};

// Type for agents field in form
export type FormAgents = FormData['agents'];

export type FormData = z.infer<typeof formSchema>;

export type ProviderHealthResponse = {
    embedding: EndpointHealth;
    embeddingProvider: string;
    provider: EndpointHealth;
    providerType: string;
};

export { agentConfigSchema, formSchema };
