#!/usr/bin/env node

const { execSync, spawn } = require('child_process')
const readline = require('readline')

const rl = readline.createInterface({
	input: process.stdin,
	output: process.stdout,
})

function question(prompt) {
	return new Promise(resolve => {
		rl.question(prompt, resolve)
	})
}

function execCommand(command) {
	try {
		return execSync(command, { encoding: 'utf8' }).trim()
	} catch (error) {
		return null
	}
}

function getGitChanges() {
	const status = execCommand('git status --porcelain')
	if (!status) {
		console.log('No changes to commit')
		return null
	}
	
	const changes = {
		added: [],
		modified: [],
		deleted: [],
		frontend: 0,
		backend: 0,
		config: 0,
		docs: 0,
	}
	
	status.split('\n').forEach(line => {
		if (!line.trim()) return
		
		const statusCode = line.substring(0, 2)
		const file = line.substring(3)
		
		switch (statusCode.trim()) {
			case 'A':
			case '??':
				changes.added.push(file)
				break
			case 'M':
				changes.modified.push(file)
				break
			case 'D':
				changes.deleted.push(file)
				break
		}
		
		// Categorize files
		if (/frontend\/|\.svelte|\.ts|\.js|\.css|\.html/.test(file)) {
			changes.frontend++
		} else if (/\.go|go\.mod|go\.sum/.test(file)) {
			changes.backend++
		} else if (/\.json|\.config|\.yml|\.yaml|prettierrc|tasks\.json|settings\.json/.test(file)) {
			changes.config++
		} else if (/README|\.md/.test(file)) {
			changes.docs++
		}
	})
	
	return changes
}

function generateCommitMessage(changes) {
	let type = 'feat'
	
	if (changes.modified.length > changes.added.length) {
		type = 'fix'
	} else if (changes.deleted.length > 0) {
		type = 'refactor'
	}
	
	const descriptions = []
	
	if (changes.frontend > 0) descriptions.push('frontend updates')
	if (changes.backend > 0) descriptions.push('backend changes')
	if (changes.config > 0) descriptions.push('configuration')
	if (changes.docs > 0) descriptions.push('documentation')
	
	const description = descriptions.length > 0 ? descriptions.join(', ') : 'project updates'
	
	return `${type}: ${description}`
}

function showChangesPreview(changes) {
	console.log('\n=== Git Changes Preview ===')
	
	if (changes.added.length > 0) {
		console.log('\x1b[32mAdded files:\x1b[0m')
		changes.added.forEach(file => console.log(`\x1b[32m  + ${file}\x1b[0m`))
	}
	
	if (changes.modified.length > 0) {
		console.log('\x1b[33mModified files:\x1b[0m')
		changes.modified.forEach(file => console.log(`\x1b[33m  ~ ${file}\x1b[0m`))
	}
	
	if (changes.deleted.length > 0) {
		console.log('\x1b[31mDeleted files:\x1b[0m')
		changes.deleted.forEach(file => console.log(`\x1b[31m  - ${file}\x1b[0m`))
	}
	
	console.log('')
}

async function main() {
	const changes = getGitChanges()
	if (!changes) {
		process.exit(1)
	}
	
	showChangesPreview(changes)
	
	let message = process.argv[2]
	
	if (!message) {
		message = generateCommitMessage(changes)
		console.log(`\x1b[36mGenerated commit message: \x1b[0m'${message}'`)
		
		const confirm = await question('Use this message? (Y/n/edit): ')
		
		if (confirm.toLowerCase() === 'n') {
			message = await question('Enter custom commit message: ')
		} else if (confirm.toLowerCase() === 'edit') {
			const customMessage = await question(`Edit message [${message}]: `)
			if (customMessage.trim()) {
				message = customMessage
			}
		}
	}
	
	console.log('\x1b[34mStaging all changes...\x1b[0m')
	execCommand('git add .')
	
	console.log(`\x1b[34mCommitting with message: '${message}'\x1b[0m`)
	execCommand(`git commit -m "${message}"`)
	
	const pushConfirm = await question('Push to remote? (Y/n): ')
	
	if (pushConfirm.toLowerCase() !== 'n') {
		console.log('\x1b[34mPushing to remote...\x1b[0m')
		execCommand('git push')
	}
	
	console.log('\x1b[32mDone! ✅\x1b[0m')
	rl.close()
}

main().catch(console.error)
