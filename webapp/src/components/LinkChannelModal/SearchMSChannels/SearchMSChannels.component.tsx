import React, {useCallback, useEffect, useState} from 'react';

import {ListItemType, MMSearch} from '@brightscout/mattermost-ui-library';

import {useDispatch} from 'react-redux';

import usePluginApi from 'hooks/usePluginApi';
import utils from 'utils';
import {debounceFunctionTimeLimit, defaultPage, defaultPerPage} from 'constants/common.constants';
import {pluginApiServiceConfigs} from 'constants/apiService.constant';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';

import {setLinkModalLoading} from 'reducers/linkModal';

import {SearchMSChannelProps} from './SearchMSChannels.types';

export const SearchMSChannels = ({setChannel, teamId}: SearchMSChannelProps) => {
    const dispatch = useDispatch();
    const {makeApiRequestWithCompletionStatus, getApiState} = usePluginApi();
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [searchChannelsPayload, setSearchChannelsPayload] = useState<SearchMSChannelsParams | null>(null);
    const [searchSuggestions, setSearchSuggestions] = useState<ListItemType[]>([]);

    useEffect(() => {
        handleClearInput();
    }, [teamId]);

    const searchChannels = ({searchFor}: {searchFor?: string}) => {
        if (searchFor && teamId) {
            const payload = {
                search: searchFor,
                page: defaultPage,
                per_page: defaultPerPage,
                teamId,
            };
            setSearchChannelsPayload(payload);
            makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.searchMSChannels.apiServiceName, payload);
            dispatch(setLinkModalLoading(true));
        }
    };

    const debouncedSearchChannels = useCallback(utils.debounce(searchChannels, debounceFunctionTimeLimit), [searchChannels]);

    const handleSearch = (val: string) => {
        if (!val) {
            setSearchSuggestions([]);
            setChannel(null);
        }
        setSearchTerm(val);
        debouncedSearchChannels({searchFor: val});
    };

    const {data: searchedChannels, isLoading: searchSuggestionsLoading} = getApiState(pluginApiServiceConfigs.searchMSChannels.apiServiceName, searchChannelsPayload as SearchMSChannelsParams);
    const handleChannelSelect = (_: any, option: ListItemType) => {
        setChannel({
            ID: option.value,
            DisplayName: option.label as string,
        });
        setSearchTerm(option.label as string);
    };

    const handleClearInput = () => {
        setSearchTerm('');
        setChannel(null);
        setSearchSuggestions([]);
    };

    useApiRequestCompletionState({
        serviceName: pluginApiServiceConfigs.searchMSChannels.apiServiceName,
        payload: searchChannelsPayload as SearchMSChannelsParams,
        handleSuccess: () => {
            if (searchedChannels) {
                const suggestions: ListItemType[] = [];
                for (const channel of searchedChannels as MSTeamsSearchResponse) {
                    suggestions.push({
                        label: channel.DisplayName,
                        value: channel.ID,
                    });
                }
                setSearchSuggestions(suggestions);
            }
            dispatch(setLinkModalLoading(false));
        },
        handleError: () => {
            dispatch(setLinkModalLoading(false));

            // TODO: Handle this error
        },
    });

    return (
        <div className='mt-24'>
            <MMSearch
                fullWidth={true}
                label='Select a channel in Microsoft Teams'
                items={searchSuggestions}
                onSelect={handleChannelSelect}
                searchValue={searchTerm}
                setSearchValue={handleSearch}
                onClearInput={handleClearInput}
                optionsLoading={searchSuggestionsLoading}
                disabled={!teamId}
            />
        </div>);
};
