import {useEffect, useRef} from 'react';

function usePreviousState(value: Record<string, string>) {
    const ref = useRef<Record<string, string>>();
    useEffect(() => {
        ref.current = value;
    }, [value]);
    return ref.current;
}

export default usePreviousState;
