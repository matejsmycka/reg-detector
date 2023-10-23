<#
Ester Niclos Ferreras
Returns nomber of web site connections
UNKNOWN - not found
OK - connections
warning - current connections greater than warning value
critical - current connection greater than critical value
#>


#
# Shell arguments
#
[CmdletBinding()]
Param(
   [Parameter(Mandatory=$True,Position=1)]
   [string]$server,
  [Parameter(Mandatory=$True,Position=2)]
   [string]$website,
   [Parameter(Mandatory=$True,Position=3)]
   [int]$warning_value,
   [Parameter(Mandatory=$True,Position=4)]
   [int]$critical_value
   )

Set-Variable OK 0 -option Constant
Set-Variable WARNING 1 -option Constant
Set-Variable CRITICAL 2 -option Constant
Set-Variable UNKNOWN 3 -option Constant


#
# ASK STATUS
#

$counter = Get-WmiObject -Class Win32_PerfFormattedData_W3SVC_WebService | Where {$_.name -match $website}

# Nagios output

$resultstring='CONNECTIONS UNKNOWN ' + $website + ' not found' 
$exit_code = $UNKNOWN
  
if ($counter -ne $null) {

	$connections=$counter.ConnectionAttemptsPersec
	
	if ($connections -gt $critical_value) {
		$status_str= 'CONNECTIONS CRITICAL '+ $website +' connections '+ $connections
		$exit_code = $CRITICAL
	}
	elseif ($connections -gt $warning_value) {
		$status_str= 'CONNECTIONS WARNING '+ $website +' connections '+ $connections
		$exit_code = $WARNING
	}
	else{
		$status_str= 'CONNECTIONS OK '+ $website +' connections '+ $connections
		$exit_code = $OK
	}
    	
	$perf_data= "ConnectionAttempts=" + $connections + ';' + $warning_value + ';' + $critical_value + "; "
	$perf_data+= "LogonAttempts=" + $counter.LogonAttemptsPersec + ';' + $warning_value + ';' + $critical_value + "; "
	$resultstring= "$status_str  |  $perf_data " 
}



Write-Host $resultstring
exit $exit_code
