#!/usr/bin/env node

const { execSync } = require("child_process")
const readline = require("readline")

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
		return execSync(command, { encoding: "utf8" }).trim()
	} catch (error) {
		return null
	}
}

async function main() {
	const status = execCommand("git status --porcelain")
	if (!status) {
		console.log("No changes to commit")
		process.exit(1)
	}

	console.log("=== Changes ===")
	console.log(execCommand("git status --short"))

	let message = await question("\nEnter commit message (or press Enter for auto-generated): ")

	if (!message) {
		const fileCount = status.split("\n").filter(line => line.trim()).length
		message = `update: ${fileCount} files changed`

		// Улучшаем сообщение на основе типов файлов
		if (/frontend\//.test(status)) {
			message = "frontend: updates and improvements"
		} else if (/\.go/.test(status)) {
			message = "core improvements"
		} else if (/\.md|README/.test(status)) {
			message = "docs: update documentation"
		} else if (/\.json|\.config/.test(status)) {
			message = "config: update project configuration"
		}
	}

	console.log(`\x1b[32mCommitting: '${message}'\x1b[0m`)

	execCommand("git add .")
	execCommand(`git commit -m "${message}"`)

	const pushConfirm = await question("Push? (Y/n): ")

	if (pushConfirm.toLowerCase() !== "n") {
		execCommand("git push")
		console.log("\x1b[32mPushed! ✅\x1b[0m")
	}

	rl.close()
}

main().catch(console.error)
