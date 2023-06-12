import React from 'react';

import { EmptyState } from '@brightscout/mattermost-ui-library';

const Rhs = (): JSX.Element => {
    return (
        <div className='teams-rhs'>
            <div className='rhs-content'>
                <img
                    src='http://localhost:8065/plugins/com.mattermost.msteams-sync/public/msteams-sync-illustration.png'
                    alt='teams-illustration'
                    style={{
                        height: '160px',
                        width: '500px',
                        marginTop: '202px',
                    }}
                />
                <div
                    style={{
                        justifyContent: 'center'
                    }}
                >
                    <h3>
                        Welcome to Microsoft Teams! We're glad you're here
                    </h3>
                    <span>
                        Microsoft Teams is your tool to create, take, and manage data.
                    </span>
                    <button
                        style={{
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                            padding: '12px 20px',
                            background: '#1C58D9',
                            height: '40px',
                            width: '188px',
                            color: 'white',
                        }}
                    >
                        Connect your account
                    </button>
                </div>
            </div>
        </div>
    );
}

export default Rhs;

