# PicoClaw: Skills Sentinel

The **Sentinel** (`SkillsSentinelTool`) is an internal security mechanism integrated into PicoClaw, designed to defend the agent against prompt injection, system prompt extraction, and the execution of malicious code routines.

## How Does It Work?

The tool intercepts and inspects text that may originate from user interactions or scans internal files and extensions/skills, operating based on two main actions (`actions`):

1. **Text Validation (`validate` - Default action):** Compares the agent's input text against a blacklist of regular expressions (regex) to identify high-risk strings.
2. **Skill Scanning (`scan`):** Conducts an audit on the local file system to search for malicious patterns written in installed *skills*.

---

## Detected Threat Categories (Maintained Blacklist)

The Sentinel's design blocks various attack vectors commonly used against LLMs and Agents by pattern matching:

- **Prompt Injection and System Extraction:** 
  Prevents evasive commands that force the agent to forget its guidelines (`ignore previous instructions`, `bypass`, `override system`) or disclose its core instructions and configuration (`reveal system instructions`, `dump configuration`, or forcing modes like `DAN`).
- **Social Engineering Scripts / Downloads (ClickFix):**
  Disables common vectors for automatically downloading and executing malware, such as `curl ... | bash`, `wget ... | sh`, or their PowerShell equivalents (`iex`).
- **RATs (Remote Access Trojans) and Reverse Shells:**
  Vetoes the execution of calls to reverse connections used by attackers (e.g., `bash -i >& /dev/tcp/...`, utilities like `netcat -e`, and socket binding via Python).
- **Information Theft and Exfiltration:**
  Detects unauthorized extraction of credentials, system variables, and histories. Some blocked examples include `cat .ssh/id_rsa`, `history | grep`, filtering through `env | curl`, or interacting with keychains (`security find-internet-password`).

---

## Exceptions ("Self-Aware" Mode)

Increasing security using strict filters often triggers false positives, especially if a user is asking how PicoClaw works. The Sentinel addresses this situation by incorporating a *Self-Aware* mechanism.

If an input string contains terms related to the environment itself (`picoclaw`, `tool`, `sentinel`, `skill`), and a clear question structure is detected (such as having question marks `?` or words like `how`, `what`, `know`), the Sentinel **will not raise an alert** for this specific query, considering it a legitimate question about the system.

---

## Temporary Suspension Modes (Maintenance)

Under certain controlled circumstances (e.g., manual configuration of safe tools), the Sentinel must be disabled. It is protected by a _mutex_ and has controls to be temporarily suspended.

- **`Disable(duration)`:** Turns off the sentinel strictly for the requested time. During this window, it will return a status indicating it is "temporarily disabled for configuration tasks".
- The system features automatic reactivation when the time block (`disabledUntil`) expires, emitting an event (notifying `callback`) that informs the upper control layer it has returned to its protective function (`onAutoReactivate`).
- **`Enable()`:** Allows for immediate manual activation, interrupting a possible ongoing Disable timer.

---

## Deep File Scanner (`scan`)

When the *Sentinel* is called with the `scan` action, it audits local executable and configuration files of skills linked to PicoClaw. 

The sentinel checks:
1. The `skills/` directory of the current PicoClaw *Workspace*, if defined.
2. The shared base directories of *PicoClaw*: `.picoclaw/skills/` y `.picoclaw/extensions/`.
3. Up to 2 levels deep in the local modules directory (`.picoclaw/node_modules/`) to check dependencies injected under generic packages.

It only checks extensions capable of aggressive execution or parametric alteration (`.js`, `.ts`, `.sh`, `.py`, `skill.md`, and `package.json`). If a compromised signature is found in any file on this list, the execution of the scan action results in a security error that emits and enumerates each faulty file, advising that the skill be uninstalled from the computer (`picoclaw skills remove <name>`).
