import { memo, useCallback, useState, useEffect, useMemo } from 'react';
import { Clock, User, Bot, Brain, Copy, ChevronDown, ChevronRight } from 'lucide-react';

import Markdown from '@/components/Markdown';
import Terminal from '@/components/Terminal';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import type { AssistantLogFragmentFragment, MessageLogFragmentFragment } from '@/graphql/types';
import { MessageLogType, ResultFormat } from '@/graphql/types';
import { cn } from '@/lib/utils';
import { formatDate } from '@/lib/utils/format';
import { copyMessageToClipboard } from '@/lib/сlipboard';

interface AcademicChatMessageProps {
    log: MessageLogFragmentFragment | AssistantLogFragmentFragment;
    searchValue?: string;
    provider?: string;
}

const containsSearchValue = (text: string | null | undefined, searchValue: string): boolean => {
    if (!text || !searchValue) return false;
    return text.toLowerCase().includes(searchValue.toLowerCase());
};

const AcademicChatMessage = ({ log, searchValue = '', provider }: AcademicChatMessageProps) => {
    const { type, createdAt, message, thinking, result, resultFormat = ResultFormat.Plain } = log;
    const isUserMessage = type === MessageLogType.Input;
    const isAssistantMessage = type === MessageLogType.Output || type === MessageLogType.Report;
    
    const [isThinkingVisible, setIsThinkingVisible] = useState(false);
    const [isDetailsVisible, setIsDetailsVisible] = useState(true);
    const [copySuccess, setCopySuccess] = useState(false);

    // Memoize search checks
    const searchChecks = useMemo(() => {
        const trimmedSearch = searchValue.trim();
        if (!trimmedSearch) {
            return { hasThinkingMatch: false, hasResultMatch: false };
        }
        
        return {
            hasThinkingMatch: containsSearchValue(thinking, trimmedSearch),
            hasResultMatch: containsSearchValue(result, trimmedSearch),
        };
    }, [searchValue, thinking, result]);

    const { hasThinkingMatch, hasResultMatch } = searchChecks;
    const shouldShowThinkingToggle = thinking && thinking.trim().length > 0;

    // Auto-expand thinking if it contains search match
    useEffect(() => {
        if (hasThinkingMatch && searchValue.trim()) {
            setIsThinkingVisible(true);
        }
    }, [hasThinkingMatch, searchValue]);

    // Auto-expand details if result contains search match
    useEffect(() => {
        if (hasResultMatch && searchValue.trim()) {
            setIsDetailsVisible(true);
        }
    }, [hasResultMatch, searchValue]);

    const toggleThinking = useCallback(() => {
        setIsThinkingVisible(prev => !prev);
    }, []);

    const toggleDetails = useCallback(() => {
        setIsDetailsVisible(prev => !prev);
    }, []);

    const handleCopy = useCallback(async () => {
        const success = await copyMessageToClipboard(log);
        if (success) {
            setCopySuccess(true);
            setTimeout(() => setCopySuccess(false), 2000);
        }
    }, [log]);

    const formatTimestamp = (timestamp: string) => {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('en-US', { 
            hour: '2-digit', 
            minute: '2-digit',
            hour12: false 
        });
    };

    const renderThinkingContent = () => {
        if (!shouldShowThinkingToggle || !isThinkingVisible) return null;

        return (
            <div className="academic-thinking-content">
                <div className="flex items-center gap-2 mb-2">
                    <Brain className="w-4 h-4 text-blue-600" />
                    <span className="text-xs font-semibold text-blue-700 uppercase tracking-wide">
                        Reasoning Process
                    </span>
                </div>
                <Markdown className="prose-sm prose-slate max-w-none" searchValue={searchValue}>
                    {thinking}
                </Markdown>
            </div>
        );
    };

    const renderDetailsContent = () => {
        if (!result || !isDetailsVisible) return null;

        return (
            <div className="mt-4">
                {resultFormat === ResultFormat.Plain && (
                    <Markdown className="prose prose-slate max-w-none leading-relaxed" searchValue={searchValue}>
                        {result}
                    </Markdown>
                )}
                {resultFormat === ResultFormat.Markdown && (
                    <Markdown className="prose prose-slate max-w-none leading-relaxed" searchValue={searchValue}>
                        {result}
                    </Markdown>
                )}
                {resultFormat === ResultFormat.Terminal && (
                    <div className="academic-code-block">
                        <Terminal
                            logs={[result as string]}
                            className="h-auto min-h-[120px] w-full bg-transparent text-slate-100"
                        />
                    </div>
                )}
            </div>
        );
    };

    if (isUserMessage) {
        return (
            <div className="academic-message-user">
                <div className="flex items-start gap-3 max-w-2xl">
                    <div className="academic-message-avatar academic-message-avatar-user">
                        <User className="w-4 h-4" />
                    </div>
                    <div className="flex-1">
                        <div className="academic-message-meta">
                            <span className="font-medium">You</span>
                            <span className="academic-timestamp">
                                {formatTimestamp(createdAt)}
                            </span>
                        </div>
                        <div className="academic-message-bubble-user">
                            {message && (
                                <Markdown className="prose prose-slate prose-invert max-w-none" searchValue={searchValue}>
                                    {message}
                                </Markdown>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    if (isAssistantMessage) {
        return (
            <div className="academic-message-assistant">
                <div className="flex items-start gap-3 max-w-3xl w-full">
                    <div className="academic-message-avatar academic-message-avatar-assistant">
                        <Bot className="w-4 h-4" />
                    </div>
                    <div className="flex-1 min-w-0">
                        <div className="academic-message-meta">
                            <span className="font-medium">Assistant</span>
                            {provider && (
                                <span className="academic-provider-badge">
                                    {provider}
                                </span>
                            )}
                            <span className="academic-timestamp">
                                {formatTimestamp(createdAt)}
                            </span>
                            <div className="ml-auto flex items-center gap-2">
                                <Tooltip>
                                    <TooltipTrigger asChild>
                                        <button
                                            onClick={handleCopy}
                                            className="p-1 rounded hover:bg-slate-100 transition-colors"
                                        >
                                            <Copy className="w-3 h-3 text-slate-400" />
                                        </button>
                                    </TooltipTrigger>
                                    <TooltipContent>
                                        {copySuccess ? 'Copied!' : 'Copy message'}
                                    </TooltipContent>
                                </Tooltip>
                            </div>
                        </div>
                        
                        <div className="academic-message-bubble-assistant">
                            {/* Thinking toggle */}
                            {shouldShowThinkingToggle && (
                                <div className="mb-3">
                                    <button
                                        onClick={toggleThinking}
                                        className="academic-thinking-toggle flex items-center gap-1"
                                    >
                                        {isThinkingVisible ? (
                                            <ChevronDown className="w-3 h-3" />
                                        ) : (
                                            <ChevronRight className="w-3 h-3" />
                                        )}
                                        {isThinkingVisible ? 'Hide reasoning' : 'Show reasoning'}
                                    </button>
                                </div>
                            )}

                            {/* Thinking content */}
                            {renderThinkingContent()}

                            {/* Main message content */}
                            {message && (
                                <div className={shouldShowThinkingToggle && isThinkingVisible ? 'mt-4' : ''}>
                                    <Markdown className="prose prose-slate max-w-none leading-relaxed" searchValue={searchValue}>
                                        {message}
                                    </Markdown>
                                </div>
                            )}

                            {/* Details content */}
                            {renderDetailsContent()}
                        </div>
                    </div>
                </div>
            </div>
        );
    }

    return null;
};

export default memo(AcademicChatMessage);
