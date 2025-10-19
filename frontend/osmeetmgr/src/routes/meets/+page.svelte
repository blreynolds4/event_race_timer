<script>
	import { onMount } from 'svelte';

	let meets = [];
	let error = null;

	onMount(async () => {
		try {
			const response = await fetch('/api/meets');
			if (!response.ok) {
				throw new Error(`Failed to fetch meets: ${response.statusText}`);
			}
			meets = await response.json();
		} catch (err) {
			error = err.message;
		}
	});
</script>

<main>
	<h1>Meets</h1>

	{#if error}
		<p class="error">{error}</p>
	{:else if meets.length === 0}
		<p>Loading meets...</p>
	{:else}
		<ul>
			{#each meets as meet}
				<li>{meet.Name}</li>
			{/each}
		</ul>
	{/if}
</main>

<style>
	.error {
		color: red;
	}
</style>
