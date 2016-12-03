import ActionTypes from './action-types';
import Immutable from 'immutable';

const initialState = Immutable.fromJS({
	ticks: []
});

export default function(state = initialState, action) {
	switch (action.type) {
		case ActionTypes.TICK:
			return state.set('ticks', state.get('ticks').push(state.payload));
		default:
			return state;
	}
}
