// import { useDispatch } from "react-redux";
import {store} from '../store/reduxStore'; 
import { showToaster, showHideLoader as toggle, showLoader as showL, hideLoader as hideL } from '../store/sidebarSlice';

// const dispatch = useDispatch();

export const showNotification = (type: 'success' | 'info' | 'error' | undefined, message: string, delay?: number) => {
    store.dispatch(showToaster({ type, message, delay }));
};
export const toggleLoader = (show?: boolean) => {
    if (show === true) {
        store.dispatch(showL());
    } else if (show === false) {
        store.dispatch(hideL());
    } else {
        store.dispatch(toggle());
    }
};

export const showLoader = () => {
    store.dispatch(showL());
};

export const hideLoader = () => {
    store.dispatch(hideL());
};
