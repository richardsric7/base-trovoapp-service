import { useState } from "react"
import { useIdleTimer } from "react-idle-timer"
/**
 * @param onIdle - function to notify user when idle timeout is close
 * @param idleTime - number of seconds to wait before user is logged out
 */
const useIdleTimeout = ({ onIdle, onPrompt, idleTime = 1 }: {onIdle: ()=> void, onPrompt: ()=> void, idleTime: number}) => {
    const idleTimeout = 1000 * idleTime;
    const [isIdle, setIdle] = useState(false)

    const idleTimer = useIdleTimer({
        timeout: idleTimeout,
        promptTimeout: idleTimeout / 3,
        onPrompt: onPrompt,
        onIdle: onIdle,
        debounce: 500
    })
    return {
        isIdle,
        setIdle,
        idleTimer
    }
}
export default useIdleTimeout;