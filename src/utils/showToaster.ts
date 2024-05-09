// import { useDispatch } from "react-redux";
import {store} from '../store/reduxStore'; 
import { showToaster, showHideLoader as toggle } from '../store/sidebarSlice';

// const dispatch = useDispatch();

export const showNotification = (type: 'success' | 'info' | 'error' | undefined, message: string, delay?: number) => {
    store.dispatch(showToaster({ type, message, delay }));
};
export const toggleLoader = () => {
    store.dispatch(toggle());
};
