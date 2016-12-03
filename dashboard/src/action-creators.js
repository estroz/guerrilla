import ActionTypes from './action-types';

export default {
	tick(payload) {
		return {
			type: ActionTypes.TICK,
			payload
		}
	}
};
