import React, {useCallback} from 'react';
import {useDispatch} from 'react-redux';

import {DialogProps} from '@brightscout/mattermost-ui-library';

import {showDialog as showDialogComponent, closeDialog} from 'reducers/dialog';

import {Dialog} from 'components';
import {DialogState} from 'types/common/store.d';

const useDialog = (dialogId: string) => {
    const dispatch = useDispatch();

    const showDialog = (props: DialogState) => dispatch(showDialogComponent({
        dialogId,
        state: props,
    }));

    const hideDialog = () => dispatch(closeDialog(dialogId));

    const DialogComponent = useCallback(({onCloseHandler, onSubmitHandler}: Pick<DialogProps, 'onCloseHandler' | 'onSubmitHandler'>) => (
        <Dialog
            id={dialogId}
            onCloseHandler={onCloseHandler}
            onSubmitHandler={onSubmitHandler}
        />
    ), [showDialog, hideDialog]);

    return {
        DialogComponent,
        showDialog,
        hideDialog,
    };
};

export default useDialog;
