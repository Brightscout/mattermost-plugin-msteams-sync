import React, {useState} from 'react';

import {LinearProgress, Modal} from '@brightscout/mattermost-ui-library';

import {useDispatch} from 'react-redux';

import usePluginApi from 'hooks/usePluginApi';

import {getCurrentTeam, getLinkModalState} from 'selectors';

import {Dialog} from 'components/Dialog';
import {pluginApiServiceConfigs} from 'constants/apiService.constant';
import useApiRequestCompletionState from 'hooks/useApiRequestCompletionState';
import {setLinkModalLoading, showLinkModal} from 'reducers/linkModal';
import useAlert from 'hooks/useAlert';

import {refetch} from 'reducers/refetchState';

import {SearchMSChannels} from './SearchMSChannels';
import {SearchMSTeams} from './SearchMSTeams';
import {SearchMMChannels} from './SearchMMChannels';

export const LinkChannelModal = ({onClose}: {onClose: () => void}) => {
    const dispatch = useDispatch();
    const showAlert = useAlert();
    const {state, makeApiRequestWithCompletionStatus} = usePluginApi();
    const {show = false, isLoading} = getLinkModalState(state);
    const {currentTeamId} = getCurrentTeam(state);

    // Show retry dialog component
    const [showRetryDialog, setShowRetryDialog] = useState(false);

    const [mMChannel, setMmChannel] = useState<Channel | null>(null);
    const [mSTeam, setMsTeam] = useState<MSTeamOrChannel | null>(null);
    const [mSChannel, setMsChannel] = useState<MSTeamOrChannel | null>(null);
    const [linkChannelsPayload, setLinkChannelsPayload] = useState<LinkChannelsPayload | null>(null);

    const handleModalClose = (preserveFields?: boolean) => {
        if (!preserveFields) {
            setMmChannel(null);
            setMsTeam(null);
            setMsChannel(null);
        }
        onClose();
    };

    const handleChannelLinking = () => {
        const payload: LinkChannelsPayload = {
            mattermostTeamID: currentTeamId || '',
            mattermostChannelID: mMChannel?.id || '',
            msTeamsTeamID: mSTeam?.ID || '',
            msTeamsChannelID: mSChannel?.ID || '',
        };
        setLinkChannelsPayload(payload);
        makeApiRequestWithCompletionStatus(pluginApiServiceConfigs.linkChannels.apiServiceName, payload);
        dispatch(setLinkModalLoading(true));
    };

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
            setShowRetryDialog(true);
        },
    });

    return (
        <>
            <Modal
                show={show}
                className='msteams-sync-modal'
                title='Link a channel'
                subtitle='Link a channel in Mattermost with a channel in Microsoft Teams.'
                primaryActionText='Link Channels'
                secondaryActionText='Cancel'
                onFooterCloseHandler={() => handleModalClose(true)}
                onHeaderCloseHandler={() => handleModalClose(true)}
                isPrimaryButtonDisabled={!mMChannel || !mSChannel || !mSTeam}
                onSubmitHandler={handleChannelLinking}
            >
                {isLoading && <LinearProgress className='fixed w-full left-0 top-100'/>}
                <SearchMMChannels
                    setChannel={setMmChannel}
                    teamId={currentTeamId}
                />
                <hr className='w-full my-32'/>
                <SearchMSTeams setMsTeam={setMsTeam}/>
                <SearchMSChannels
                    setChannel={setMsChannel}
                    teamId={mSTeam?.ID}
                />
            </Modal>
            <Dialog
                show={showRetryDialog}
                destructive={true}
                primaryButtonText='Try Again'
                secondaryButtonText='Cancel'
                title='Unable to link channels'
                onSubmitHandler={() => {
                    setShowRetryDialog(false);
                    dispatch(showLinkModal());
                }}
                onCloseHandler={() => setShowRetryDialog(false)}
            >
                {'We were not able to link the selected channels. Please try again.'}
            </Dialog>
        </>
    );
};
