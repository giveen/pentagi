import { AlertCircle, Cpu, Loader2, Play, Save, Trash2, Wand2 } from 'lucide-react';

import ConfirmationDialog from '@/components/shared/confirmation-dialog';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { StatusCard } from '@/components/ui/status-card';
import { ReasoningEffort } from '@/graphql/types';
import { cn } from '@/lib/utils';

import {
    FormComboboxItem,
    FormInputNumberItem,
    FormInputStringItem,
    FormModelComboboxItem,
} from './provider-form-fields';
import TestResultsDialog from './provider-test-dialog';
import { useProviderForm } from './use-provider-form';

const SettingsProvider = () => {
    const {
        agentTypes,
        agentTypesMap,
        availableModels,
        control,
        currentAgentKey,
        currentTuningAgentKey,
        endpointHealth,
        error,
        form,
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
        providerName,
        providers,
        selectedType,
        setIsDeleteDialogOpen,
        setIsTestDialogOpen,
        setValue,
        testResults,
    } = useProviderForm();

    if (loading) {
        return (
            <StatusCard
                description="Please wait while we fetch provider configuration"
                icon={<Loader2 className="text-muted-foreground size-16 animate-spin" />}
                title="Loading provider data..."
            />
        );
    }

    if (error) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="size-4" />
                <AlertTitle>Error loading provider data</AlertTitle>
                <AlertDescription>{error.message}</AlertDescription>
            </Alert>
        );
    }

    return (
        <>
            <div className="flex flex-col gap-4">
                <div className="flex flex-col gap-2">
                    <h2 className="flex items-center gap-2 text-lg font-semibold">
                        <Cpu className="text-muted-foreground size-5" />
                        {isNew ? 'New Provider' : 'Provider Settings'}
                    </h2>

                    <div className="text-muted-foreground">
                        {isNew
                            ? 'Configure a new language model provider'
                            : 'Update provider settings and configuration'}
                    </div>
                </div>

                <Form {...form}>
                    <form
                        className="flex flex-col gap-6"
                        id="provider-form"
                        onSubmit={handleFormSubmit(handleSubmit)}
                    >
                        {/* Error Alert */}
                        {mutationError && (
                            <Alert variant="destructive">
                                <AlertCircle className="size-4" />
                                <AlertTitle>Error</AlertTitle>
                                <AlertDescription>
                                    {mutationError instanceof Error ? (
                                        mutationError.message
                                    ) : (
                                        <div className="whitespace-pre-line">{mutationError}</div>
                                    )}
                                </AlertDescription>
                            </Alert>
                        )}

                        {/* Form fields */}
                        <FormComboboxItem
                            allowCustom={false}
                            control={control}
                            description="The type of language model provider"
                            disabled={isLoading || !!selectedType}
                            label="Type"
                            name="type"
                            options={providers}
                            placeholder="Select provider"
                        />

                        <FormInputStringItem
                            control={control}
                            description="A unique name for your provider configuration"
                            disabled={isLoading}
                            label="Name"
                            name="name"
                            placeholder="Enter provider name"
                        />

                        <FormInputStringItem
                            control={control}
                            description="Optional. Override API base URL for this provider in dashboard settings"
                            disabled={isLoading}
                            label="API URL"
                            name="apiUrl"
                            placeholder="https://api.example.com/v1"
                        />

                        <FormInputStringItem
                            control={control}
                            description="Optional. Override API key for this provider in dashboard settings"
                            disabled={isLoading}
                            label="API Key"
                            name="apiKey"
                            placeholder="Enter API key"
                        />

                        <div className="flex flex-col gap-3">
                            <Button
                                disabled={isLoading || isTestLoading || isAgentTestLoading || isEndpointHealthLoading}
                                onClick={handleCheckEndpoint}
                                type="button"
                                variant="outline"
                            >
                                {isEndpointHealthLoading ? (
                                    <Loader2 className="size-4 animate-spin" />
                                ) : (
                                    <Cpu className="size-4" />
                                )}
                                {isEndpointHealthLoading ? 'Checking Endpoint...' : 'Check Provider Endpoint'}
                            </Button>

                            {endpointHealth && (
                                <Alert variant="default">
                                    <AlertCircle className="size-4" />
                                    <AlertTitle>Endpoint Health</AlertTitle>
                                    <AlertDescription>
                                        <div className="space-y-2 text-sm">
                                            <div>
                                                <span className="font-medium">Provider ({endpointHealth.providerType}): </span>
                                                {endpointHealth.provider.reachable ? 'reachable' : 'unreachable'}
                                                {endpointHealth.provider.url ? ` at ${endpointHealth.provider.url}` : ''}
                                                {endpointHealth.provider.models !== undefined
                                                    ? `, models: ${endpointHealth.provider.models}`
                                                    : ''}
                                                {endpointHealth.provider.statusCode !== undefined
                                                    ? `, status: ${endpointHealth.provider.statusCode}`
                                                    : ''}
                                                {endpointHealth.provider.error ? ` (${endpointHealth.provider.error})` : ''}
                                            </div>
                                            <div>
                                                <span className="font-medium">
                                                    Embedding ({endpointHealth.embeddingProvider}):
                                                </span>{' '}
                                                {endpointHealth.embedding.reachable ? 'reachable' : 'unreachable'}
                                                {endpointHealth.embedding.model
                                                    ? `, model: ${endpointHealth.embedding.model}`
                                                    : ''}
                                                {endpointHealth.embedding.url ? ` at ${endpointHealth.embedding.url}` : ''}
                                                {endpointHealth.embedding.models !== undefined
                                                    ? `, models: ${endpointHealth.embedding.models}`
                                                    : ''}
                                                {endpointHealth.embedding.statusCode !== undefined
                                                    ? `, status: ${endpointHealth.embedding.statusCode}`
                                                    : ''}
                                                {endpointHealth.embedding.error ? ` (${endpointHealth.embedding.error})` : ''}
                                            </div>
                                        </div>
                                    </AlertDescription>
                                </Alert>
                            )}
                        </div>

                        {/* Agents Configuration Section */}
                        <div className="flex flex-col gap-4">
                            <div>
                                <h3 className="text-lg font-medium">Agent Configurations</h3>
                                <p className="text-muted-foreground text-sm">Configure settings for each agent type</p>
                            </div>

                            <Accordion
                                className="w-full"
                                type="multiple"
                            >
                                {agentTypes.map((agentKey) => (
                                    <AccordionItem
                                        key={agentKey}
                                        value={agentKey}
                                    >
                                        <AccordionTrigger className="group text-left hover:no-underline">
                                            <div className="flex w-full items-center justify-between gap-2">
                                                <span className="group-hover:underline">
                                                    {agentKey === 'embedding'
                                                        ? 'Embedder'
                                                        : agentKey
                                                              .replaceAll(/([A-Z])/g, ' $1')
                                                              .replace(/^./, (s) => s.toUpperCase())}
                                                </span>
                                                {agentTypesMap[agentKey] && (
                                                    <span
                                                        className={cn(
                                                            'hover:bg-accent hover:text-accent-foreground mr-2 flex items-center gap-1 rounded border px-2 py-1 text-xs',
                                                            currentTuningAgentKey !== null &&
                                                                'pointer-events-none cursor-not-allowed opacity-50',
                                                        )}
                                                        onClick={(event) => {
                                                            if (currentTuningAgentKey !== null) {
                                                                return;
                                                            }

                                                            event.stopPropagation();
                                                            handleAutoTune(agentKey);
                                                        }}
                                                    >
                                                        {currentTuningAgentKey === agentKey ? (
                                                            <Loader2 className="size-4 animate-spin" />
                                                        ) : (
                                                            <Wand2 className="size-4" />
                                                        )}
                                                        <span className="no-underline! hover:no-underline!">
                                                            {currentTuningAgentKey === agentKey
                                                                ? 'Tuning...'
                                                                : 'Auto Tune'}
                                                        </span>
                                                    </span>
                                                )}
                                                {agentTypesMap[agentKey] && (
                                                    <span
                                                        className={cn(
                                                            'hover:bg-accent hover:text-accent-foreground mr-2 flex items-center gap-1 rounded border px-2 py-1 text-xs',
                                                            (isTestLoading || isAgentTestLoading) &&
                                                                'pointer-events-none cursor-not-allowed opacity-50',
                                                        )}
                                                        onClick={(event) => {
                                                            if (isTestLoading || isAgentTestLoading) {
                                                                return;
                                                            }

                                                            event.stopPropagation();
                                                            handleTestAgent(agentKey);
                                                        }}
                                                    >
                                                        {isAgentTestLoading && currentAgentKey === agentKey ? (
                                                            <Loader2 className="size-4 animate-spin" />
                                                        ) : (
                                                            <Play className="size-4" />
                                                        )}
                                                        <span className="no-underline! hover:no-underline!">
                                                            {isAgentTestLoading && currentAgentKey === agentKey
                                                                ? 'Testing...'
                                                                : 'Test'}
                                                        </span>
                                                    </span>
                                                )}
                                            </div>
                                        </AccordionTrigger>
                                        <AccordionContent className="flex flex-col gap-4 pt-4">
                                            <div className="grid grid-cols-1 gap-4 p-px md:grid-cols-2">
                                                {/* Model field */}
                                                <FormModelComboboxItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Model"
                                                    name={`agents.${agentKey}.model`}
                                                    onOptionSelect={(option) => {
                                                        const price = option?.price;

                                                        setValue(
                                                            `agents.${agentKey}.price.input` as const,
                                                            price?.input ?? null,
                                                        );
                                                        setValue(
                                                            `agents.${agentKey}.price.output` as const,
                                                            price?.output ?? null,
                                                        );
                                                        setValue(
                                                            `agents.${agentKey}.price.cacheRead` as const,
                                                            price?.cacheRead ?? null,
                                                        );
                                                        setValue(
                                                            `agents.${agentKey}.price.cacheWrite` as const,
                                                            price?.cacheWrite ?? null,
                                                        );
                                                    }}
                                                    options={availableModels}
                                                    placeholder="Select or enter model name"
                                                />

                                                {/* Temperature field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Temperature"
                                                    max="2"
                                                    min="0"
                                                    name={`agents.${agentKey}.temperature`}
                                                    placeholder="0.7"
                                                    step="0.1"
                                                />

                                                {/* Max Tokens field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Max Tokens"
                                                    min="1"
                                                    name={`agents.${agentKey}.maxTokens`}
                                                    placeholder="1000"
                                                    valueType="integer"
                                                />

                                                {/* Top P field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Top P"
                                                    max="1"
                                                    min="0"
                                                    name={`agents.${agentKey}.topP`}
                                                    placeholder="0.9"
                                                    step="0.01"
                                                />

                                                {/* Top K field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Top K"
                                                    min="1"
                                                    name={`agents.${agentKey}.topK`}
                                                    placeholder="40"
                                                    valueType="integer"
                                                />

                                                {/* Min Length field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Min Length"
                                                    min="0"
                                                    name={`agents.${agentKey}.minLength`}
                                                    placeholder="0"
                                                    valueType="integer"
                                                />

                                                {/* Max Length field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Max Length"
                                                    min="1"
                                                    name={`agents.${agentKey}.maxLength`}
                                                    placeholder="2000"
                                                    valueType="integer"
                                                />

                                                {/* Repetition Penalty field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Repetition Penalty"
                                                    max="2"
                                                    min="0"
                                                    name={`agents.${agentKey}.repetitionPenalty`}
                                                    placeholder="1.0"
                                                    step="0.01"
                                                />

                                                {/* Frequency Penalty field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Frequency Penalty"
                                                    max="2"
                                                    min="0"
                                                    name={`agents.${agentKey}.frequencyPenalty`}
                                                    placeholder="0.0"
                                                    step="0.01"
                                                />

                                                {/* Presence Penalty field */}
                                                <FormInputNumberItem
                                                    control={control}
                                                    disabled={isLoading}
                                                    label="Presence Penalty"
                                                    max="2"
                                                    min="0"
                                                    name={`agents.${agentKey}.presencePenalty`}
                                                    placeholder="0.0"
                                                    step="0.01"
                                                />
                                            </div>

                                            {/* Reasoning Configuration */}
                                            <div className="col-span-full p-px">
                                                <div className="mt-6 flex flex-col gap-4">
                                                    <h4 className="text-sm font-medium">Reasoning Configuration</h4>
                                                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                                                        {/* Reasoning Effort field */}
                                                        <FormField
                                                            control={control}
                                                            name={`agents.${agentKey}.reasoning.effort`}
                                                            render={({ field }) => (
                                                                <FormItem>
                                                                    <FormLabel>Reasoning Effort</FormLabel>
                                                                    <Select
                                                                        defaultValue={field.value ?? 'none'}
                                                                        disabled={isLoading}
                                                                        onValueChange={(value) =>
                                                                            field.onChange(
                                                                                value !== 'none' ? value : null,
                                                                            )
                                                                        }
                                                                    >
                                                                        <FormControl>
                                                                            <SelectTrigger>
                                                                                <SelectValue placeholder="Select effort level (optional)" />
                                                                            </SelectTrigger>
                                                                        </FormControl>
                                                                        <SelectContent>
                                                                            <SelectItem value="none">
                                                                                Not selected
                                                                            </SelectItem>
                                                                            <SelectItem value={ReasoningEffort.Low}>
                                                                                Low
                                                                            </SelectItem>
                                                                            <SelectItem value={ReasoningEffort.Medium}>
                                                                                Medium
                                                                            </SelectItem>
                                                                            <SelectItem value={ReasoningEffort.High}>
                                                                                High
                                                                            </SelectItem>
                                                                        </SelectContent>
                                                                    </Select>
                                                                    <FormMessage />
                                                                </FormItem>
                                                            )}
                                                        />

                                                        {/* Reasoning Max Tokens field */}
                                                        <FormInputNumberItem
                                                            control={control}
                                                            disabled={isLoading}
                                                            label="Reasoning Max Tokens"
                                                            min="1"
                                                            name={`agents.${agentKey}.reasoning.maxTokens`}
                                                            placeholder="1000"
                                                            valueType="integer"
                                                        />
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Price Configuration */}
                                            <div className="col-span-full p-px">
                                                <div className="mt-6 flex flex-col gap-4">
                                                    <h4 className="text-sm font-medium">Price Configuration</h4>
                                                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                                                        {/* Price Input field */}
                                                        <FormInputNumberItem
                                                            control={control}
                                                            description="Price per 1M input tokens"
                                                            disabled={isLoading}
                                                            label="Input Price"
                                                            min="0"
                                                            name={`agents.${agentKey}.price.input`}
                                                            placeholder="0.001"
                                                            step="0.000001"
                                                        />

                                                        {/* Price Output field */}
                                                        <FormInputNumberItem
                                                            control={control}
                                                            description="Price per 1M output tokens"
                                                            disabled={isLoading}
                                                            label="Output Price"
                                                            min="0"
                                                            name={`agents.${agentKey}.price.output`}
                                                            placeholder="0.002"
                                                            step="0.000001"
                                                        />

                                                        {/* Cache Read Price field */}
                                                        <FormInputNumberItem
                                                            control={control}
                                                            description="Price per 1M cached read tokens"
                                                            disabled={isLoading}
                                                            label="Cache Read Price"
                                                            min="0"
                                                            name={`agents.${agentKey}.price.cacheRead`}
                                                            placeholder="0.0001"
                                                            step="0.000001"
                                                        />

                                                        {/* Cache Write Price field */}
                                                        <FormInputNumberItem
                                                            control={control}
                                                            description="Price per 1M cache write tokens"
                                                            disabled={isLoading}
                                                            label="Cache Write Price"
                                                            min="0"
                                                            name={`agents.${agentKey}.price.cacheWrite`}
                                                            placeholder="0.00015"
                                                            step="0.000001"
                                                        />
                                                    </div>
                                                </div>
                                            </div>
                                        </AccordionContent>
                                    </AccordionItem>
                                ))}
                            </Accordion>
                        </div>
                    </form>
                </Form>
            </div>

            {/* Sticky buttons at bottom */}
            <div className="bg-background sticky -bottom-4 -mx-4 mt-4 -mb-4 flex items-center border-t p-4 shadow-lg">
                <div className="flex gap-2">
                    {/* Delete button - only show when editing existing provider */}
                    {!isNew && (
                        <Button
                            disabled={isLoading}
                            onClick={handleDelete}
                            type="button"
                            variant="destructive"
                        >
                            {isDeleteLoading ? (
                                <Loader2 className="size-4 animate-spin" />
                            ) : (
                                <Trash2 className="size-4" />
                            )}
                            {isDeleteLoading ? 'Deleting...' : 'Delete'}
                        </Button>
                    )}
                    <Button
                        disabled={isLoading || isTestLoading || isAgentTestLoading}
                        onClick={() => handleTest()}
                        type="button"
                        variant="outline"
                    >
                        {isTestLoading ? <Loader2 className="size-4 animate-spin" /> : <Play className="size-4" />}
                        {isTestLoading ? 'Testing...' : 'Test'}
                    </Button>
                </div>

                <div className="ml-auto flex gap-2">
                    <Button
                        disabled={isLoading}
                        onClick={handleBack}
                        type="button"
                        variant="outline"
                    >
                        Cancel
                    </Button>
                    <Button
                        disabled={isLoading}
                        form="provider-form"
                        type="submit"
                        variant="secondary"
                    >
                        {isLoading ? <Loader2 className="size-4 animate-spin" /> : <Save className="size-4" />}
                        {isLoading ? 'Saving...' : isNew ? 'Create Provider' : 'Update Provider'}
                    </Button>
                </div>
            </div>

            <TestResultsDialog
                handleOpenChange={setIsTestDialogOpen}
                isOpen={isTestDialogOpen}
                results={testResults}
            />

            <ConfirmationDialog
                cancelText="Cancel"
                confirmText="Delete"
                handleConfirm={handleConfirmDelete}
                handleOpenChange={setIsDeleteDialogOpen}
                isOpen={isDeleteDialogOpen}
                itemName={providerName}
                itemType="provider"
            />

            <ConfirmationDialog
                cancelText="Stay"
                confirmIcon={undefined}
                confirmText="Leave"
                confirmVariant="destructive"
                description="You have unsaved changes. Are you sure you want to leave without saving?"
                handleConfirm={handleConfirmLeave}
                handleOpenChange={handleLeaveDialogOpenChange}
                isOpen={isLeaveDialogOpen}
                title="Discard changes?"
            />
        </>
    );
};

export default SettingsProvider;
