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
    const [mmChannel, setMMChannel] = useState<MMTeamOrChannel | null>(null);
    const [msTeam, setMSTeam] = useState<MSTeamOrChannel | null>(null);
    const [msChannel, setMSChannel] = useState<MSTeamOrChannel | null>(null);

    const handleModalClose = () => {
        setMMChannel(null);
        setMSTeam(null);
        setMSChannel(null);
        onClose();
    };

    return (
        <Modal
            show={show}
            title='Link a channel'
            subtitle='Link a channel in Mattermost with a channel in Microsoft Teams'
            primaryActionText='Link Channels'
            secondaryActionText='Cancel'
            onFooterCloseHandler={handleModalClose}
            onHeaderCloseHandler={handleModalClose}
            isPrimaryButtonDisabled={!mmChannel || !msChannel || !msTeam}
            onSubmitHandler={() => {
                // TODO: handle channel linking
            }}
        >
            {isLoading && <LinearProgress className='fixed w-full left-0 top-100'/>}
            <SearchMMChannels
                setChannel={setMMChannel}
                teamId={currentTeam}
            />
            <hr className='w-full my-32'/>
            <SearchMSTeams setMSTeam={setMSTeam}/>
            <SearchMSChannels
                setChannel={setMSChannel}
                teamId={msTeam?.ID}
            />
        </Modal>
    );
};
