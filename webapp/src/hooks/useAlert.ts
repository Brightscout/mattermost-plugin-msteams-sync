import {useDispatch} from 'react-redux';

import {alertSeverity} from 'constants/common.constants';
import {showAlert as showAlertAction} from 'reducers/snackbar';
import {SnackbarColor} from 'components/Snackbar/Snackbar.types';
import {IconName} from 'components';

const useAlert = () => {
    const dispatch = useDispatch();

    /**
	 * Show snackbar on RHs
	 * @param payload Alert message and severity
	 */
    const showAlert = ({
        message,
        severity = alertSeverity.default,
        icon = 'tick',
    }: {
        severity?: SnackbarColor;
        message: string;
        icon?: IconName;
    }) => {
        dispatch(
            showAlertAction({
                message,
                severity,
                icon,
            }),
        );
    };

    return showAlert;
};

export default useAlert;
