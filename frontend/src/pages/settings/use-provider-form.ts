import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useForm, useFormState, useWatch } from 'react-hook-form';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';

import type { AgentConfigInput, AgentsConfigInput } from '@/graphql/types';
import {
    AgentConfigType,
    useAutoTuneParamsLazyQuery,
    useCreateProviderMutation,
    useDeleteProviderMutation,
    useSettingsProvidersQuery,
    useTestAgentMutation,
    useTestProviderMutation,
    useUpdateProviderMutation,
} from '@/graphql/types';

import type { FormAgents, FormData } from './provider-form-schema';
import { formSchema } from './provider-form-schema';
import {
    agentTypesMap,
    extractAgentTypes,
    fallbackAgentTypes,
    normalizeGraphQLData,
    normalizeProviderType,
    transformFormToGraphQL,
} from './provider-helpers';
import type { Provider } from './provider-helpers';
import { useProviderHealth } from './use-provider-health';

export function useProviderForm() {
    const { providerId } = useParams<{ providerId: string }>();
    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();
    const { data, error, loading } = useSettingsProvidersQuery();
    const [createProvider, { error: createError, loading: isCreateLoading }] = useCreateProviderMutation();
    const [updateProvider, { error: updateError, loading: isUpdateLoading }] = useUpdateProviderMutation();
    const [deleteProvider, { error: deleteError, loading: isDeleteLoading }] = useDeleteProviderMutation();
    const [testProvider, { error: testError, loading: isTestLoading }] = useTestProviderMutation();
    const [testAgent, { error: agentTestError, loading: isAgentTestLoading }] = useTestAgentMutation();
    const [getAutoTuneParams] = useAutoTuneParamsLazyQuery();

    const { checkEndpoint, endpointHealth, isEndpointHealthLoading } = useProviderHealth();

    const [currentAgentKey, setCurrentAgentKey] = useState<null | string>(null);
    const [currentTuningAgentKey, setCurrentTuningAgentKey] = useState<null | string>(null);
    const [submitError, setSubmitError] = useState<null | string>(null);
    const [isTestDialogOpen, setIsTestDialogOpen] = useState(false);
    const [testResults, setTestResults] = useState<any>(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isLeaveDialogOpen, setIsLeaveDialogOpen] = useState(false);
    const [pendingBrowserBack, setPendingBrowserBack] = useState(false);
    const allowBrowserLeaveRef = useRef(false);
    const hasPushedBlockerStateRef = useRef(false);

    const isNew = providerId === 'new';
    const isLoading = isCreateLoading || isUpdateLoading || isDeleteLoading;

    const form = useForm<FormData>({
        defaultValues: {
            agents: {},
            apiKey: undefined,
            apiUrl: undefined,
            name: undefined,
            type: undefined,
        },
        resolver: zodResolver(formSchema),
    });

    const { control, formState, handleSubmit: handleFormSubmit, reset, setValue, trigger, watch } = form;

    const { isDirty } = useFormState({ control });

    // Maintain a blocker state at the top of history when form is dirty
    useEffect(() => {
        if (isDirty && !hasPushedBlockerStateRef.current) {
            window.history.pushState({ __pentagiBlock__: true }, '');
            hasPushedBlockerStateRef.current = true;
        }
    }, [isDirty]);

    // Intercept browser back using popstate when form is dirty
    useEffect(() => {
        const handlePopState = () => {
            if (!isDirty) {
                return;
            }

            if (allowBrowserLeaveRef.current) {
                // Allow single leave without blocking
                allowBrowserLeaveRef.current = false;

                return;
            }

            // User navigated back off the blocker entry; go forward to stay
            setPendingBrowserBack(true);
            setIsLeaveDialogOpen(true);
            // Return to the blocker entry
            window.history.forward();
        };

        window.addEventListener('popstate', handlePopState, { capture: true });

        return () => {
            window.removeEventListener('popstate', handlePopState, { capture: true });
        };
    }, [isDirty]);

    // Watch selected type
    const selectedType = useWatch({ control, name: 'type' });

    // Watch provider name for delete confirmation dialog
    const providerName = useWatch({ control, name: 'name' });

    // Read query parameters for form initialization (stable)
    const formQueryParams = useMemo(
        () => ({
            id: searchParams.get('id'),
            type: searchParams.get('type'),
        }),
        [searchParams],
    );

    // Get dynamic agent types from data
    const getAgentTypes = () => {
        // Try to get agents from specific sources in priority order
        const agentsSource =
            // For new providers, use default provider for selected type
            (isNew &&
                selectedType &&
                data?.settingsProviders?.default?.[selectedType as keyof typeof data.settingsProviders.default]
                    ?.agents) ||
            // For existing providers, use current provider's agents
            (!isNew &&
                providerId &&
                data?.settingsProviders?.userDefined?.find((p: Provider) => p.id == providerId)?.agents) ||
            // Fallback to any available default provider
            (data?.settingsProviders?.default &&
                Object.values(data.settingsProviders.default).find((provider) => provider?.agents)?.agents) ||
            null;

        // Extract and return agent types, or fallback to hardcoded list
        return extractAgentTypes(agentsSource) ?? fallbackAgentTypes;
    };

    const agentTypes = getAgentTypes();

    // Get available models filtered by selected provider type
    const availableModels = useMemo(() => {
        if (!data?.settingsProviders?.models || !selectedType) {
            return [];
        }

        // Filter models by selected provider type
        const { models } = data.settingsProviders;
        const providerModels = models[selectedType as keyof typeof models];

        if (!providerModels?.length) {
            return [];
        }

        return providerModels
            .map((model: any) => ({
                name: model.name,
                price: model.price
                    ? {
                          cacheRead: model.price.cacheRead ?? 0,
                          cacheWrite: model.price.cacheWrite ?? 0,
                          input: model.price.input ?? 0,
                          output: model.price.output ?? 0,
                      }
                    : null,
                thinking: model.thinking,
            }))
            .filter((model) => model.name)
            .sort((a, b) => a.name.localeCompare(b.name));
    }, [data, selectedType]);

    // Fill agents when provider type is selected (only for new providers)
    useEffect(() => {
        if (!isNew || !selectedType || !data?.settingsProviders?.default || availableModels.length === 0) {
            return;
        }

        const defaultProvider =
            data.settingsProviders.default[selectedType as keyof typeof data.settingsProviders.default];

        if (defaultProvider?.agents) {
            const agents = Object.fromEntries(
                Object.entries(defaultProvider.agents)
                    .filter(([key]) => key !== '__typename')
                    .map(([key, data]) => {
                        const agent = { ...data };

                        // Check if the model from defaultProvider exists in availableModels
                        if (agent.model && !availableModels.find((m) => m.name === agent.model)) {
                            // Use first available model if default model not found
                            agent.model = availableModels[0]?.name || agent.model;
                        }

                        return [key, agent];
                    }),
            );

            setValue('agents', normalizeGraphQLData(agents) as FormAgents);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [availableModels, data, isNew, selectedType]);

    // Update query parameter when type changes (only for new providers)
    useEffect(() => {
        if (!isNew) {
            // Clear query parameters for existing providers
            if (searchParams.size > 0) {
                setSearchParams({});
            }

            return;
        }

        // Don't update query params if we're copying from existing provider
        const queryId = searchParams.get('id');

        if (queryId) {
            return;
        }

        // Don't update query params on initial load if we're reading from query params
        const queryType = normalizeProviderType(searchParams.get('type'));

        if (!selectedType && queryType) {
            return;
        }

        // Update query parameter based on selected type
        setSearchParams((prev) => {
            const params = new URLSearchParams(prev);

            if (selectedType) {
                params.set('type', selectedType);
            } else {
                params.delete('type');
            }

            return params;
        });
    }, [selectedType, setSearchParams, isNew, searchParams]);

    // Fill form with data when available
    useEffect(() => {
        if (!data?.settingsProviders) {
            return;
        }

        const providers = data.settingsProviders;

        if (isNew || !providerId) {
            // For new provider, start with empty form but check for type query parameter
            const queryType = normalizeProviderType(formQueryParams.type);
            const queryId = formQueryParams.id;

            // If we have an id in query params, copy from existing provider
            if (queryId && data?.settingsProviders?.userDefined) {
                const sourceProvider = data.settingsProviders.userDefined.find((p: Provider) => p.id == queryId);

                if (sourceProvider) {
                    const { agents, name, type: sourceType } = sourceProvider;

                    reset({
                        agents: agents ? (normalizeGraphQLData(agents) as FormAgents) : {},
                        apiKey: sourceProvider.apiKey ?? undefined,
                        apiUrl: sourceProvider.apiUrl ?? undefined,
                        name: `${name} (Copy)`,
                        type: sourceType ?? undefined,
                    });

                    return;
                }
            } else if (queryType && data?.settingsProviders?.default) {
                const defaultProvider =
                    data.settingsProviders.default[queryType as keyof typeof data.settingsProviders.default];

                reset({
                    agents: defaultProvider?.agents ? (normalizeGraphQLData(defaultProvider.agents) as FormAgents) : {},
                    apiKey: defaultProvider?.apiKey ?? undefined,
                    apiUrl: defaultProvider?.apiUrl ?? undefined,
                    name: undefined,
                    type: queryType ?? undefined,
                });
            }

            // Default new provider form - but only if selectedType is not set
            // to avoid conflicts with agent filling useEffect
            if (!selectedType) {
                reset({
                    agents: {},
                    apiKey: undefined,
                    apiUrl: undefined,
                    name: undefined,
                    type: queryType ?? undefined,
                });
            }

            return;
        }

        const provider = providers.userDefined?.find((provider: Provider) => provider.id == providerId);

        if (!provider) {
            navigate('/settings/providers');

            return;
        }

        const { agents, apiKey, apiUrl, name, type } = provider;

        reset({
            agents: agents ? (normalizeGraphQLData(agents) as FormAgents) : {},
            apiKey: apiKey || undefined,
            apiUrl: apiUrl || undefined,
            name: name || undefined,
            type: type || undefined,
        });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [data, formQueryParams, isNew, providerId, selectedType]);

    const handleCheckEndpoint = async () => {
        const selectedProviderType = normalizeProviderType(watch('type'));
        if (!selectedProviderType) {
            setSubmitError('Provider type is required for endpoint health check.');

            return;
        }

        try {
            setSubmitError(null);
            await checkEndpoint({
                apiKey: watch('apiKey') || undefined,
                apiUrl: watch('apiUrl') || undefined,
                embeddingModel: watch('agents.embedding.model') || undefined,
                type: selectedProviderType,
            });
        } catch (error) {
            console.error('Endpoint health check error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while checking endpoint health');
        }
    };

    const handleSubmit = async () => {
        // Note: getValues() excludes disabled fields, watch() includes them
        const formData = watch();

        try {
            setSubmitError(null);

            const mutationData = transformFormToGraphQL(formData);

            if (isNew) {
                await createProvider({
                    refetchQueries: ['settingsProviders'],
                    variables: mutationData,
                });
            } else {
                await updateProvider({
                    refetchQueries: ['settingsProviders'],
                    variables: {
                        ...mutationData,
                        providerId: providerId!,
                    },
                });
            }

            navigate('/settings/providers');
        } catch (error) {
            console.error('Submit error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while saving');
        }
    };

    const handleDelete = () => {
        if (isNew || !providerId) {
            return;
        }

        setIsDeleteDialogOpen(true);
    };

    const handleConfirmDelete = async () => {
        if (isNew || !providerId) {
            return;
        }

        try {
            setSubmitError(null);
            await deleteProvider({
                refetchQueries: ['settingsProviders'],
                variables: { providerId },
            });
            navigate('/settings/providers');
        } catch (error) {
            console.error('Delete error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while deleting');
        }
    };

    // Build a human-readable validation error message from react-hook-form errors
    const buildValidationError = (errors: any): string => {
        const formatFieldName = (fieldPath: string): string =>
            fieldPath
                .split('.')
                .map((part) => part.charAt(0).toUpperCase() + part.slice(1).replaceAll(/([A-Z])/g, ' $1'))
                .join(' → ');

        const errorMessages = Object.entries(errors)
            .map(([field, error]: [string, any]) => {
                if (error?.message) {
                    return `• ${formatFieldName(field)}: ${error.message}`;
                }

                if (error && typeof error === 'object') {
                    return Object.entries(error)
                        .map(([subField, subError]: [string, any]) => {
                            if (subError?.message) {
                                return `• ${formatFieldName(`${field}.${subField}`)}: ${subError.message}`;
                            }

                            if (subError && typeof subError === 'object') {
                                return Object.entries(subError)
                                    .map(([nestedField, nestedError]: [string, any]) => {
                                        if (nestedError?.message) {
                                            return `• ${formatFieldName(`${field}.${subField}.${nestedField}`)}: ${nestedError.message}`;
                                        }

                                        return null;
                                    })
                                    .filter(Boolean)
                                    .join('\n');
                            }

                            return null;
                        })
                        .filter(Boolean)
                        .join('\n');
                }

                return null;
            })
            .filter(Boolean)
            .join('\n');

        return `Please fix the following validation errors:\n\n${errorMessages}`;
    };

    const handleTest = async () => {
        const isValid = await trigger();

        if (!isValid) {
            setSubmitError(buildValidationError(formState.errors));

            return;
        }

        try {
            setSubmitError(null);

            // Get form data and transform it - including disabled fields
            const formData = watch();
            const { agents, apiKey, apiUrl, type } = transformFormToGraphQL(formData);
            const result = await testProvider({
                variables: {
                    agents,
                    apiKey,
                    apiUrl,
                    type,
                },
            });

            setTestResults(result.data?.testProvider);
            setIsTestDialogOpen(true);
        } catch (error) {
            console.error('Test error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while testing');
        }
    };

    // Test a single agent (uses testAgent where supported, otherwise falls back to filtered provider test)
    const handleTestAgent = async (agentKey: string) => {
        const isValid = await trigger();

        if (!isValid) {
            setSubmitError(buildValidationError(formState.errors));

            return;
        }

        try {
            setSubmitError(null);
            setCurrentAgentKey(agentKey);
            // Note: getValues() excludes disabled fields, watch() includes them
            const formData = watch();
            const { agents, type } = transformFormToGraphQL(formData);

            const agent = agents[agentKey as keyof AgentsConfigInput] as AgentConfigInput;

            const singleResult = await testAgent({
                variables: { agent, agentType: agentTypesMap[agentKey] ?? AgentConfigType.Simple, type },
            });
            setTestResults({ [agentKey]: singleResult.data?.testAgent });
            setIsTestDialogOpen(true);
            setCurrentAgentKey(null);

            return;
        } catch (error) {
            console.error('Test error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while testing');
            setCurrentAgentKey(null);
        }
    };

    const handleAutoTune = async (agentKey: string) => {
        const agentType = agentTypesMap[agentKey];
        const formData = watch();
        const modelName = formData.agents?.[agentKey]?.model ?? '';

        if (!agentType || !modelName) {
            setSubmitError('Please set a model for this agent before auto-tuning.');

            return;
        }

        try {
            setSubmitError(null);
            setCurrentTuningAgentKey(agentKey);
            const result = await getAutoTuneParams({ variables: { agentType, modelName } });

            if (result.data?.autoTuneParams) {
                const p = result.data.autoTuneParams;
                setValue(`agents.${agentKey}.temperature` as const, p.temperature, { shouldDirty: true });
                setValue(`agents.${agentKey}.topP` as const, p.topP, { shouldDirty: true });
                setValue(`agents.${agentKey}.topK` as const, p.topK, { shouldDirty: true });
                setValue(`agents.${agentKey}.frequencyPenalty` as const, p.frequencyPenalty, { shouldDirty: true });
                setValue(`agents.${agentKey}.presencePenalty` as const, p.presencePenalty, { shouldDirty: true });
                setValue(`agents.${agentKey}.repetitionPenalty` as const, p.repetitionPenalty, { shouldDirty: true });
            }
        } catch (error) {
            console.error('Auto-tune error:', error);
            setSubmitError(error instanceof Error ? error.message : 'An error occurred while auto-tuning');
        } finally {
            setCurrentTuningAgentKey(null);
        }
    };

    const handleBack = () => {
        if (isDirty) {
            setIsLeaveDialogOpen(true);

            return;
        }

        navigate('/settings/providers');
    };

    const handleConfirmLeave = () => {
        if (pendingBrowserBack) {
            allowBrowserLeaveRef.current = true;
            setPendingBrowserBack(false);
            // Skip the blocker entry and go to the real previous page
            window.history.go(-2);

            return;
        }

        navigate('/settings/providers');
    };

    const handleLeaveDialogOpenChange = (open: boolean) => {
        if (!open && pendingBrowserBack) {
            setPendingBrowserBack(false);
        }

        setIsLeaveDialogOpen(open);
    };

    const mutationError = createError || updateError || deleteError || testError || agentTestError || submitError;

    const providers = data?.settingsProviders?.models
        ? Object.keys(data.settingsProviders.models).filter((key) => key !== '__typename')
        : [];

    return {
        agentTypes,
        agentTypesMap,
        availableModels,
        control,
        currentAgentKey,
        currentTuningAgentKey,
        data,
        endpointHealth,
        error,
        form,
        formState,
        handleAutoTune,
        handleBack,
        handleCheckEndpoint,
        handleConfirmDelete,
        handleConfirmLeave,
        handleDelete,
        handleFormSubmit,
        handleLeaveDialogOpenChange,
        handleSubmit,
        handleTest,
        handleTestAgent,
        isAgentTestLoading,
        isDeleteDialogOpen,
        isDeleteLoading,
        isEndpointHealthLoading,
        isLeaveDialogOpen,
        isLoading,
        isNew,
        isTestDialogOpen,
        isTestLoading,
        loading,
        mutationError,
        providerId,
        providerName,
        providers,
        selectedType,
        setIsDeleteDialogOpen,
        setIsTestDialogOpen,
        setValue,
        testResults,
    };
}
