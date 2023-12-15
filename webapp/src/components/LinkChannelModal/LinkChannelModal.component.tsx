import React, {useEffect, useState} from 'react';

import {LinearProgress, Modal} from '@brightscout/mattermost-ui-library';

import {useDispatch} from 'react-redux';

import usePluginApi from 'hooks/usePluginApi';

import {getCurrentTeam, getLinkModalState} from 'selectors';

import {pluginApiServiceConfigs} from 'constants/apiService.constant';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import {hideLinkModal, preserveState, resetState, setLinkModalLoading, showLinkModal} from 'reducers/linkModal';
import useAlert from 'hooks/useAlert';

import {refetch} from 'reducers/refetchState';

import useDialog from 'hooks/useDialog';

import {DialogsIds} from 'constants/common.constants';

import {SearchMSChannels} from './SearchMSChannels';
import {SearchMSTeams} from './SearchMSTeams';
import {SearchMMChannels} from './SearchMMChannels';

export const LinkChannelModal = () => {
    const dispatch = useDispatch();
    const showAlert = useAlert();
    const {state, makeApiRequestWithCompletionStatus} = usePluginApi();
    const {show = false, isLoading} = getLinkModalState(state);
    const {currentTeamId} = getCurrentTeam(state);

    const [mmChannel, setMMChannel] = useState<MMTeamOrChannel | null>(null);
    const [msTeam, setMSTeam] = useState<MSTeamOrChannel | null>(null);
    const [msChannel, setMSChannel] = useState<MSTeamOrChannel | null>(null);
    const [linkChannelsPayload, setLinkChannelsPayload] = useState<LinkChannelsPayload | null>(null);

    const handleModalClose = (preserve?: boolean) => {
        if (!preserve) {
            setMMChannel(null);
            setMSTeam(null);
            setMSChannel(null);
        }
        dispatch(resetState());
        dispatch(hideLinkModal());
    };

    const handleChannelLinking = () => {
        const payload: LinkChannelsPayload = {
            mattermostTeamID: currentTeamId || '',
            mattermostChannelID: mmChannel?.id || '',
            msTeamsTeamID: msTeam?.ID || '',
            msTeamsChannelID: msChannel?.ID || '',
        };
        setLinkChannelsPayload(payload);
        makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.linkChannels.apiServiceName, payload);
        dispatch(setLinkModalLoading(true));
    };

    const {DialogComponent, showDialog, hideDialog} = useDialog(DialogsIds.retryLink);

    useApiRequestCompletionState({
        serviceName: pluginApiServiceConfigs.linkChannels.apiServiceName,
        payload: linkChannelsPayload as LinkChannelsPayload,
        handleSuccess: () => {
            dispatch(setLinkModalLoading(false));
            handleModalClose();
            dispatch(refetch());
            showAlert({
                message: 'Successfully linked channels',
                severity: 'success',
            });
        },
        handleError: () => {
            dispatch(setLinkModalLoading(false));
            handleModalClose(true);
            showDialog({
                title: 'Unable to link channels',
                description: 'We were not able to link the selected channels. Please try again.',
                primaryButtonText: 'Try Again',
                secondaryButtonText: 'Cancel',
            });
        },
    });

    return (
        <>
            <Modal
                show={show}
                className='msteams-sync-modal'
                title='Link a channel'
                subtitle='Link a channel in Mattermost with a channel in Microsoft Teams'
                primaryActionText='Link Channels'
                secondaryActionText='Cancel'
                onFooterCloseHandler={handleModalClose}
                onHeaderCloseHandler={handleModalClose}
                isPrimaryButtonDisabled={!mmChannel || !msChannel || !msTeam}
                onSubmitHandler={handleChannelLinking}
                backdrop={true}
            >
                {isLoading && <LinearProgress className='fixed w-full left-0 top-100'/>}
                <SearchMMChannels
                    setChannel={setMMChannel}
                    teamId={currentTeamId}
                />
                <hr className='w-full my-32'/>
                <SearchMSTeams setMSTeam={setMSTeam}/>
                <SearchMSChannels
                    setChannel={setMSChannel}
                    teamId={msTeam?.ID}
                />
            </Modal>
            <DialogComponent
                onSubmitHandler={() => {
                    dispatch(preserveState({
                        mmChannel: mmChannel?.displayName ?? '',
                        msChannel: msChannel?.DisplayName ?? '',
                        msTeam: msTeam?.DisplayName ?? '',
                    }));
                    dispatch(showLinkModal());
                    hideDialog();
                }}
                onCloseHandler={() => {
                    hideDialog();
                    dispatch(resetState());
                }}
            />
        </>
    );
};
