!macro NSIS_HOOK_PREINSTALL
  nsExec::ExecToLog 'taskkill /F /IM pit-desktop.exe'
  nsExec::ExecToLog 'taskkill /F /IM pit.exe'
  nsExec::ExecToLog 'taskkill /F /IM PIT.exe'
  Sleep 1500
!macroend

!macro NSIS_HOOK_PREUNINSTALL
  nsExec::ExecToLog 'taskkill /F /IM pit-desktop.exe'
  nsExec::ExecToLog 'taskkill /F /IM pit.exe'
  Sleep 800
!macroend
