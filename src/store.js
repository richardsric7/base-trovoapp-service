import { configureStore } from '@reduxjs/toolkit';
import sidebarSlice from './reducers/sidebarSlice';

export default configureStore({
  reducer: {
    sidebar: sidebarSlice,
  },
});
