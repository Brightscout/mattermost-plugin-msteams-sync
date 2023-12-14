import React from 'react';

import {Button, Icon as UILibIcon, Tooltip} from '@brightscout/mattermost-ui-library';

import {General as MMConstants} from 'mattermost-redux/constants';

import {Icon} from 'components/Icon';

import {LinkedChannelCardProps} from './LinkedChannelCard.types';

import './LinkedChannelCard.styles.scss';

export const LinkedChannelCard = ({msTeamsChannelName, msTeamsTeamName, mattermostChannelName, mattermostTeamName, mattermostChannelType}: LinkedChannelCardProps) => (
    <div className='px-16 py-12 border-t-1 d-flex gap-4 msteams-linked-channel'>
        <div className='msteams-linked-channel__link-icon d-flex align-items-center flex-column justify-center'>
            <Icon iconName='link'/>
        </div>
        <div className='d-flex flex-column gap-6 msteams-linked-channel__body'>
            <div className='d-flex gap-8 align-items-center'>
                {mattermostChannelType === MMConstants.PRIVATE_CHANNEL ? <Icon iconName='lock'/> : <Icon iconName='globe'/>}
                <Tooltip
                    placement='left'
                    text={mattermostChannelName}
                >
                    <h5 className='my-0 msteams-linked-channel__body-values'>{mattermostChannelName}</h5>
                </Tooltip>
                <Tooltip
                    placement='left'
                    text={mattermostTeamName}
                >
                    <h5 className='my-0 opacity-6 msteams-linked-channel__body-values'>{mattermostTeamName}</h5>
                </Tooltip>
            </div>
            <div className='d-flex gap-8 align-items-center'>
                <Icon iconName='msTeams'/>
                <Tooltip
                    placement='left'
                    text={msTeamsChannelName}
                >
                    <h5 className='my-0 msteams-linked-channel__body-values'>{msTeamsChannelName}</h5>
                </Tooltip>
                <Tooltip
                    placement='left'
                    text={msTeamsTeamName}
                >
                    <h5 className='my-0 opacity-6 msteams-linked-channel__body-values'>{msTeamsTeamName}</h5>
                </Tooltip>
            </div>
        </div>
        <Button
            variant='text'
            aria-label='unlink channel'
            className='msteams-linked-channel__unlink-icon'
        >
            <UILibIcon
                name='Unlink'
                size={16}
            />
        </Button>
    </div>
);
