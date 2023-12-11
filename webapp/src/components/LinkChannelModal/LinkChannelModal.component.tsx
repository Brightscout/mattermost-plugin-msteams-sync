import React, {useEffect, useState} from 'react';

import {LinearProgress, Modal} from '@brightscout/mattermost-ui-library';

import usePluginApi from 'hooks/usePluginApi';

import {getCurrentTeam, getLinkModalState} from 'selectors';

import {SearchMMChannels} from './SearchMMChannels';
import {SearchMSTeams} from './SearchMSTeams';
import {SearchMSChannels} from './SearchMSChannels';

export const LinkChannelModal = ({onClose}: {onClose: () => void}) => {
    const {state} = usePluginApi();
    const {show = false, isLoading} = getLinkModalState(state);
    const currentTeam = getCurrentTeam(state);
    const [mMChannel, setMmChannel] = useState<Channel | null>(null);
    const [mSTeam, setMsTeam] = useState<MSTeamOrChannel | null>(null);
    const [mSChannel, setMsChannel] = useState<MSTeamOrChannel | null>(null);

    const handleModalClose = () => {
        setMmChannel(null);
        setMsTeam(null);
        setMsChannel(null);
        onClose();
    };

    return (
        <Modal
            show={show}
            title='Link a channel'
            subtitle='Link a channel in Mattermost with a channel in Microsoft Teams.'
            primaryActionText='Link Channels'
            secondaryActionText='Cancel'
            onFooterCloseHandler={handleModalClose}
            onHeaderCloseHandler={handleModalClose}
            isPrimaryButtonDisabled={!mMChannel || !mSChannel || !mSTeam}
            onSubmitHandler={() => {
                // TODO: handle channel linking
            }}
        >
            {isLoading && <LinearProgress className='fixed w-full left-0 top-100'/>}
            <SearchMMChannels
                setChannel={setMmChannel}
                teamId={currentTeam}
            />
            <hr className='w-full my-32'/>
            <SearchMSTeams setMsTeam={setMsTeam}/>
            <SearchMSChannels
                setChannel={setMsChannel}
                teamId={mSTeam?.ID}
            />
        </Modal>
    );
};
