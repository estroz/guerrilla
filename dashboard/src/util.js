// ...x -> {...x: x}
export function Actions(...types) {
	const actions = {};
	for (let i = 0; i < types.length; i++) {
		actions[types[i]] = types[i];
	}
	return actions;
}
