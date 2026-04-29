import { CheckCircle, Clock, XCircle } from 'lucide-react';

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';

interface TestResultsDialogProps {
    handleOpenChange: (isOpen: boolean) => void;
    isOpen: boolean;
    results: any;
}

const getStatusIcon = (result: boolean) => {
    if (result === true) {
        return <CheckCircle className="size-4 text-green-500" />;
    } else if (result === false) {
        return <XCircle className="size-4 text-red-500" />;
    } else {
        return <Clock className="size-4 text-yellow-500" />;
    }
};

const getStatusColor = (result: boolean) => {
    if (result === true) {
        return 'text-green-600';
    } else if (result === false) {
        return 'text-red-600';
    } else {
        return 'text-yellow-600';
    }
};

// Component to render test results dialog
const TestResultsDialog = ({ handleOpenChange, isOpen, results }: TestResultsDialogProps) => {
    if (!results) {
        return null;
    }

    const agentResults = Object.entries(results)
        .filter(([key]) => key !== '__typename')
        .map(([agentType, agentData]: [string, any]) => ({
            agentType,
            tests: agentData?.tests || [],
        }));

    return (
        <Dialog
            onOpenChange={handleOpenChange}
            open={isOpen}
        >
            <DialogContent className="flex max-h-[80vh] max-w-4xl flex-col">
                <DialogHeader className="shrink-0">
                    <DialogTitle>Provider Test Results</DialogTitle>
                </DialogHeader>
                <div className="flex flex-1 flex-col gap-6 overflow-y-auto">
                    <Accordion
                        className="w-full"
                        type="multiple"
                    >
                        {agentResults.map(({ agentType, tests }) => {
                            const testsCount = tests.length;
                            const successTestsCount = tests.filter((test: any) => test.result === true).length;

                            return (
                                <AccordionItem
                                    key={agentType}
                                    value={agentType}
                                >
                                    <AccordionTrigger className="text-left">
                                        <div className="mr-4 flex w-full items-center justify-between">
                                            <span className="text-lg font-semibold capitalize">{agentType}</span>
                                            <span className="text-muted-foreground text-sm">
                                                {successTestsCount}/{testsCount} tests passed
                                            </span>
                                        </div>
                                    </AccordionTrigger>
                                    <AccordionContent>
                                        <div className="flex flex-col gap-3 pt-2">
                                            {tests.map((test: any, index: number) => (
                                                <div
                                                    className="rounded-lg border p-3"
                                                    key={index}
                                                >
                                                    <div className="mb-2 flex items-start justify-between">
                                                        <div className="flex items-center gap-2">
                                                            {getStatusIcon(test.result)}
                                                            <span className="font-medium">{test.name}</span>
                                                            {test.type && (
                                                                <span className="text-muted-foreground text-sm">
                                                                    ({test.type})
                                                                </span>
                                                            )}
                                                        </div>
                                                        <div className="text-muted-foreground flex items-center gap-3 text-sm">
                                                            {test.reasoning !== undefined && (
                                                                <span>Reasoning: {test.reasoning ? 'Yes' : 'No'}</span>
                                                            )}
                                                            {test.streaming !== undefined && (
                                                                <span>Streaming: {test.streaming ? 'Yes' : 'No'}</span>
                                                            )}
                                                            {test.latency && <span>Latency: {test.latency}ms</span>}
                                                        </div>
                                                    </div>
                                                    <div
                                                        className={`text-sm font-medium ${getStatusColor(test.result)}`}
                                                    >
                                                        Result:{' '}
                                                        {test.result === true
                                                            ? 'Success'
                                                            : test.result === false
                                                              ? 'Failed'
                                                              : 'Unknown'}
                                                    </div>
                                                    {test.error && (
                                                        <div className="mt-2 rounded border border-red-200 bg-red-50 p-2 text-sm text-red-700">
                                                            <strong>Error:</strong> {test.error}
                                                        </div>
                                                    )}
                                                </div>
                                            ))}
                                            {tests.length === 0 && (
                                                <div className="text-muted-foreground py-4 text-center">
                                                    No tests available for this agent
                                                </div>
                                            )}
                                        </div>
                                    </AccordionContent>
                                </AccordionItem>
                            );
                        })}
                    </Accordion>
                </div>
            </DialogContent>
        </Dialog>
    );
};

export default TestResultsDialog;
