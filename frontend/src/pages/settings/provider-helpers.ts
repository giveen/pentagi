import type { AgentConfigInput, AgentsConfigInput, ProviderConfigFragmentFragment } from '@/graphql/types';
import { AgentConfigType, ProviderType, ReasoningEffort } from '@/graphql/types';

import type { FormData } from './provider-form-schema';

export type Provider = ProviderConfigFragmentFragment;

export const validProviderTypes = new Set<string>(Object.values(ProviderType));

export const normalizeProviderType = (value: null | string | undefined): ProviderType | undefined => {
    if (!value || !validProviderTypes.has(value)) {
        return undefined;
    }

    return value as ProviderType;
};

// Convert camelCase key to display name (e.g., 'simpleJson' -> 'Simple Json')
export const getName = (key: string): string =>
    key.replaceAll(/([A-Z])/g, ' $1').replace(/^./, (item) => item.toUpperCase());

export const getAgentDisplayName = (key: string): string => {
    if (key === 'embedding') {
        return 'Embedder';
    }

    return getName(key);
};

// Helper function to convert string to ReasoningEffort enum
export const getReasoningEffort = (effort: null | string | undefined): null | ReasoningEffort => {
    if (!effort) {
        return null;
    }

    switch (effort.toLowerCase()) {
        case 'high': {
            return ReasoningEffort.High;
        }

        case 'low': {
            return ReasoningEffort.Low;
        }

        case 'medium': {
            return ReasoningEffort.Medium;
        }

        default: {
            return null;
        }
    }
};

// Helper function to convert form data to GraphQL input
export const transformFormToGraphQL = (
    formData: FormData,
): {
    agents: AgentsConfigInput;
    apiKey?: string;
    apiUrl?: string;
    name: string;
    type: ProviderType;
} => {
    const agents = Object.entries(formData.agents || {})
        .filter(([key, data]) => key !== '__typename' && data?.model)
        .reduce((configs, [key, data]) => {
            const config: AgentConfigInput = {
                frequencyPenalty: data?.frequencyPenalty ?? null,
                maxLength: data?.maxLength ?? null,
                maxTokens: data?.maxTokens ?? null,
                minLength: data?.minLength ?? null,
                model: data!.model,
                presencePenalty: data?.presencePenalty ?? null,
                price:
                    data?.price &&
                    typeof data?.price.input === 'number' &&
                    typeof data?.price.output === 'number' &&
                    typeof data?.price.cacheRead === 'number' &&
                    typeof data?.price.cacheWrite === 'number'
                        ? {
                              cacheRead: data.price.cacheRead,
                              cacheWrite: data.price.cacheWrite,
                              input: data.price.input,
                              output: data.price.output,
                          }
                        : null,
                reasoning: data?.reasoning
                    ? {
                          effort: getReasoningEffort(data?.reasoning.effort),
                          maxTokens: data?.reasoning.maxTokens ?? null,
                      }
                    : null,
                repetitionPenalty: data?.repetitionPenalty ?? null,
                temperature: data?.temperature ?? null,
                topK: data?.topK ?? null,
                topP: data?.topP ?? null,
            };

            return { ...configs, [key]: config };
        }, {} as AgentsConfigInput);

    const providerType = normalizeProviderType(formData.type);
    if (!providerType) {
        throw new Error('Invalid provider type selected');
    }

    return {
        agents,
        apiKey: formData.apiKey || undefined,
        apiUrl: formData.apiUrl || undefined,
        name: formData.name,
        type: providerType,
    };
};

// Helper function to recursively remove __typename from objects
export const normalizeGraphQLData = (obj: unknown): unknown => {
    if (obj === null || obj === undefined) {
        return obj;
    }

    if (Array.isArray(obj)) {
        return obj.map(normalizeGraphQLData);
    }

    if (typeof obj === 'object') {
        return Object.fromEntries(
            Object.entries(obj)
                .filter(([key]) => key !== '__typename')
                .map(([key, value]) => [key, normalizeGraphQLData(value)]),
        );
    }

    return obj;
};

// Static mapping of agent keys to GraphQL enum types
export const agentTypesMap: Record<string, AgentConfigType> = {
    adviser: AgentConfigType.Adviser,
    assistant: AgentConfigType.Assistant,
    coder: AgentConfigType.Coder,
    enricher: AgentConfigType.Enricher,
    generator: AgentConfigType.Generator,
    installer: AgentConfigType.Installer,
    pentester: AgentConfigType.Pentester,
    primaryAgent: AgentConfigType.PrimaryAgent,
    refiner: AgentConfigType.Refiner,
    reflector: AgentConfigType.Reflector,
    searcher: AgentConfigType.Searcher,
    simple: AgentConfigType.Simple,
    simpleJson: AgentConfigType.SimpleJson,
};

export const fallbackAgentTypes: string[] = ['embedding', ...Object.keys(agentTypesMap)];

// Helper function to extract agent types from agents object
export const extractAgentTypes = (agents: unknown): null | string[] => {
    if (!agents || typeof agents !== 'object') {
        return null;
    }

    const types = Object.entries(agents)
        .filter(([key, data]) => key !== '__typename' && data)
        .map(([key]) => key)
        .sort();

    if (!types.includes('embedding')) {
        types.unshift('embedding');
    }

    return types.length > 0 ? types : null;
};
