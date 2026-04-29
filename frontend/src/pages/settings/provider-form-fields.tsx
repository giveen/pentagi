import { Check, ChevronsUpDown, Lightbulb } from 'lucide-react';
import { useState } from 'react';
import { useController } from 'react-hook-form';

import { Button } from '@/components/ui/button';
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command';
import { FormControl, FormDescription, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';

// Universal field components using useController
interface ControllerProps {
    control: any;
    disabled?: boolean;
    name: string;
}

interface BaseFieldProps extends ControllerProps {
    label: string;
}

interface BaseInputProps {
    placeholder?: string;
}

interface NumberInputProps extends BaseInputProps {
    max?: string;
    min?: string;
    step?: string;
}

export interface FormInputNumberItemProps extends BaseFieldProps, NumberInputProps {
    description?: string;
    valueType?: 'float' | 'integer';
}

export interface FormInputStringItemProps extends BaseFieldProps, BaseInputProps {
    description?: string;
}

export interface FormComboboxItemProps extends BaseFieldProps, BaseInputProps {
    allowCustom?: boolean;
    contentClass?: string;
    description?: string;
    options: string[];
}

export interface ModelOption {
    name: string;
    price?: null | { cacheRead: number; cacheWrite: number; input: number; output: number };
    thinking?: boolean;
}

export interface FormModelComboboxItemProps extends BaseFieldProps, BaseInputProps {
    allowCustom?: boolean;
    contentClass?: string;
    description?: string;
    onOptionSelect?: (option: ModelOption) => void;
    options: ModelOption[];
}

export const FormInputStringItem: React.FC<FormInputStringItemProps> = ({
    control,
    description,
    disabled,
    label,
    name,
    placeholder,
}) => {
    const { field, fieldState } = useController({
        control,
        defaultValue: undefined,
        disabled,
        name,
    });

    return (
        <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
                <Input
                    {...field}
                    placeholder={placeholder}
                    value={field.value ?? ''}
                />
            </FormControl>
            {description && <FormDescription>{description}</FormDescription>}
            {fieldState.error && <FormMessage>{fieldState.error.message}</FormMessage>}
        </FormItem>
    );
};

export const FormInputNumberItem: React.FC<FormInputNumberItemProps> = ({
    control,
    description,
    disabled,
    label,
    max,
    min,
    name,
    placeholder,
    step,
    valueType = 'float',
}) => {
    const { field, fieldState } = useController({
        control,
        defaultValue: undefined,
        disabled,
        name,
    });

    const parseValue = (value: string) => {
        if (value === '') {
            return null;
        }

        return valueType === 'float' ? Number.parseFloat(value) : Number.parseInt(value);
    };

    return (
        <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
                <Input
                    {...field}
                    max={max}
                    min={min}
                    onChange={(event) => {
                        const { value } = event.target;
                        field.onChange(parseValue(value));
                    }}
                    placeholder={placeholder}
                    step={step}
                    type="number"
                    value={field.value ?? ''}
                />
            </FormControl>
            {description && <FormDescription>{description}</FormDescription>}
            {fieldState.error && <FormMessage>{fieldState.error.message}</FormMessage>}
        </FormItem>
    );
};

export const FormComboboxItem: React.FC<FormComboboxItemProps> = ({
    allowCustom = true,
    contentClass,
    control,
    description,
    disabled,
    label,
    name,
    options,
    placeholder,
}) => {
    const { field, fieldState } = useController({
        control,
        defaultValue: undefined,
        disabled,
        name,
    });

    const [isOpen, setIsOpen] = useState(false);
    const [search, setSearch] = useState('');

    const filteredOptions = options.filter((option) => option?.toLowerCase().includes(search?.toLowerCase()));
    const displayValue = field.value ?? '';

    return (
        <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
                <Popover
                    onOpenChange={setIsOpen}
                    open={isOpen}
                >
                    <PopoverTrigger asChild>
                        <Button
                            className={cn('w-full justify-between', !displayValue && 'text-muted-foreground')}
                            disabled={disabled}
                            variant="outline"
                        >
                            {displayValue || placeholder}
                            <ChevronsUpDown className="opacity-50" />
                        </Button>
                    </PopoverTrigger>
                    <PopoverContent
                        align="start"
                        className={cn(contentClass, 'p-0')}
                        style={{
                            maxHeight: 'var(--radix-popover-content-available-height)',
                            width: 'var(--radix-popover-trigger-width)',
                        }}
                    >
                        <Command>
                            <CommandInput
                                className="h-9"
                                onValueChange={setSearch}
                                placeholder={`Search ${label.toLowerCase()}...`}
                                value={search}
                            />
                            <CommandList>
                                <CommandEmpty>
                                    <div className="py-2 text-center">
                                        <p className="text-muted-foreground text-sm">No {label.toLowerCase()} found.</p>
                                        {search && allowCustom && (
                                            <Button
                                                className="mt-2"
                                                onClick={() => {
                                                    field.onChange(search);
                                                    setIsOpen(false);
                                                    setSearch('');
                                                }}
                                                size="sm"
                                                variant="ghost"
                                            >
                                                Use "{search}" as custom {label.toLowerCase()}
                                            </Button>
                                        )}
                                    </div>
                                </CommandEmpty>
                                <CommandGroup>
                                    {filteredOptions.map((option) => (
                                        <CommandItem
                                            key={option}
                                            onSelect={() => {
                                                field.onChange(option);
                                                setIsOpen(false);
                                                setSearch('');
                                            }}
                                            value={option}
                                        >
                                            {option}
                                            <Check
                                                className={cn(
                                                    'ml-auto',
                                                    displayValue === option ? 'opacity-100' : 'opacity-0',
                                                )}
                                            />
                                        </CommandItem>
                                    ))}
                                </CommandGroup>
                            </CommandList>
                        </Command>
                    </PopoverContent>
                </Popover>
            </FormControl>
            {description && <FormDescription>{description}</FormDescription>}
            {fieldState.error && <FormMessage>{fieldState.error.message}</FormMessage>}
        </FormItem>
    );
};

export const FormModelComboboxItem: React.FC<FormModelComboboxItemProps> = ({
    allowCustom = true,
    contentClass,
    control,
    description,
    disabled,
    label,
    name,
    onOptionSelect,
    options,
    placeholder,
}) => {
    const { field, fieldState } = useController({
        control,
        defaultValue: undefined,
        disabled,
        name,
    });

    const [isOpen, setIsOpen] = useState(false);
    const [search, setSearch] = useState('');

    const filteredOptions = options.filter((option) => option.name?.toLowerCase().includes(search?.toLowerCase()));
    const displayValue = field.value ?? '';

    const formatPrice = (price?: null | { cacheRead: number; cacheWrite: number; input: number; output: number }): string => {
        if (!price || ((!price.input || price.input === 0) && (!price.output || price.output === 0))) {
            return 'free';
        }

        const formatValue = (value: number): string => value.toFixed(6).replace(/\.?0+$/, '');
        const basePrice = `$${formatValue(price.input)}/$${formatValue(price.output)}`;
        const hasCachePrices = (price.cacheRead && price.cacheRead > 0) || (price.cacheWrite && price.cacheWrite > 0);

        if (hasCachePrices) {
            const cacheParts: string[] = [];

            if (price.cacheRead && price.cacheRead > 0) {
                cacheParts.push(`R:$${formatValue(price.cacheRead)}`);
            }

            if (price.cacheWrite && price.cacheWrite > 0) {
                cacheParts.push(`W:$${formatValue(price.cacheWrite)}`);
            }

            return `${basePrice} (${cacheParts.join(', ')})`;
        }

        return basePrice;
    };

    return (
        <FormItem>
            <FormLabel>{label}</FormLabel>
            <FormControl>
                <Popover
                    onOpenChange={setIsOpen}
                    open={isOpen}
                >
                    <div className="flex w-full">
                        <Input
                            className="rounded-r-none border-r-0 focus-visible:z-10"
                            disabled={disabled}
                            onChange={(event) => field.onChange(event.target.value)}
                            placeholder={placeholder}
                            value={displayValue}
                        />
                        <PopoverTrigger asChild>
                            <Button
                                className="rounded-l-none border-l-0 px-3 hover:z-10"
                                disabled={disabled}
                                type="button"
                                variant="outline"
                            >
                                <ChevronsUpDown className="size-4 opacity-50" />
                            </Button>
                        </PopoverTrigger>
                        <PopoverContent
                            align="end"
                            className={cn(contentClass, 'w-80 p-0 sm:w-[480px] md:w-[640px]')}
                        >
                            <Command>
                                <CommandInput
                                    className="h-9"
                                    onValueChange={setSearch}
                                    placeholder={`Search ${label.toLowerCase()}...`}
                                    value={search}
                                />
                                <CommandList>
                                    <CommandEmpty>
                                        <div className="py-2 text-center">
                                            <p className="text-muted-foreground text-sm">
                                                No {label.toLowerCase()} found.
                                            </p>
                                            {search && allowCustom && (
                                                <Button
                                                    className="mt-2"
                                                    onClick={() => {
                                                        field.onChange(search);
                                                        setIsOpen(false);
                                                        setSearch('');
                                                    }}
                                                    size="sm"
                                                    variant="ghost"
                                                >
                                                    Use "{search}" as custom {label.toLowerCase()}
                                                </Button>
                                            )}
                                        </div>
                                    </CommandEmpty>
                                    <CommandGroup>
                                        {filteredOptions.map((option) => (
                                            <CommandItem
                                                key={option.name}
                                                onSelect={() => {
                                                    field.onChange(option.name);
                                                    onOptionSelect?.(option);
                                                    setIsOpen(false);
                                                    setSearch('');
                                                }}
                                                value={option.name}
                                            >
                                                <div className="flex w-full min-w-0 items-center justify-between gap-2">
                                                    <div className="flex min-w-0 items-center gap-2">
                                                        <span className="truncate">{option.name}</span>
                                                        {option.thinking && (
                                                            <Lightbulb className="text-muted-foreground size-3" />
                                                        )}
                                                    </div>
                                                    <span className="text-muted-foreground shrink-0 text-xs whitespace-nowrap">
                                                        {formatPrice(option.price)}
                                                    </span>
                                                </div>
                                                <Check
                                                    className={cn(
                                                        'ml-auto',
                                                        displayValue === option.name ? 'opacity-100' : 'opacity-0',
                                                    )}
                                                />
                                            </CommandItem>
                                        ))}
                                    </CommandGroup>
                                </CommandList>
                            </Command>
                        </PopoverContent>
                    </div>
                </Popover>
            </FormControl>
            {description && <FormDescription>{description}</FormDescription>}
            {fieldState.error && <FormMessage>{fieldState.error.message}</FormMessage>}
        </FormItem>
    );
};
