package utils

/*
#include <stdio.h>
#include <windows.h>
#include <Windows.h>

int FreeDll(HMODULE dllModule)
{
	if (dllModule != NULL)
	{
		FreeLibrary(dllModule);
		//FreeLibraryAndExitThread(dllModule,0);
	}
	return 1;
}

HMODULE LoadDllModule(char* dllPath, char* dllName)
{
	//TCHAR path[MAX_PATH] = {dllPath};
	SetDllDirectory(dllPath);
    HMODULE dllModule = LoadLibraryA(dllName);
	return dllModule;
}

FARPROC loadDllFunc(HMODULE dllModule, char* dllFunc)
{
    FARPROC entry = GetProcAddress(dllModule, dllFunc);
	return entry;
}

// Read EPProjResult
typedef VARIANT (WINAPI* EPProjResultFunc)(char* jobDir, char* runnumber, char* Product, char* Key, char* VarName, char* resdate,int ResultsType);
double readRes(FARPROC entry,char* ws, char* runNo, char* prod, char* sp, char* varName, char* date)
{
	int num = 0;
	EPProjResultFunc read = (EPProjResultFunc) entry;
	VARIANT x = read(ws, runNo, prod, sp, varName, date, num);
	return x.dblVal;
}

// Read EPStochasticSummary
typedef VARIANT (WINAPI* EPStoSummaryFunc)(char* jobDir, char* runnumber, char* Product, char* Key, char* VarName, char* resdate, char* SummaryName);
double readStoSummary(FARPROC entry,char* ws, char* runNo, char* prod, char* sp, char* varName, char* date, char* summaryName)
{
	EPStoSummaryFunc read = (EPStoSummaryFunc) entry;
	VARIANT x = read(ws, runNo, prod, sp, varName, date, summaryName);
	return x.dblVal;
}


*/
import "C"

import (
	"unsafe"
)

type DllModule C.HMODULE
type DllFunc C.FARPROC

func LoadDllModule(dllPath, dllName string) DllModule {
	dPath := C.CString(dllPath)
	dName := C.CString(dllName)
	dllModule := DllModule(C.LoadDllModule(dPath, dName))
	C.free(unsafe.Pointer(dPath))
	C.free(unsafe.Pointer(dName))
	return dllModule
}

func FreeDll(dllModule DllModule) {
	_ = int(C.FreeDll(C.HMODULE(dllModule)))
}

func LoadDllFunc(dllModule DllModule, dllFunc string) DllFunc {
	dFunc := C.CString(dllFunc)
	entry := DllFunc(C.loadDllFunc(C.HMODULE(dllModule), dFunc))
	C.free(unsafe.Pointer(dFunc))
	return entry
}

func ProjResultFloat(entry DllFunc, ws, run, prod, sp, varName, calDate string) float64 {
	workSpace := C.CString(ws)
	runNo := C.CString(run)
	prodName := C.CString(prod)
	spCode := C.CString(sp)
	varStr := C.CString(varName)
	calD := C.CString(calDate)

	resDb := C.readRes(entry, workSpace, runNo, prodName, spCode, varStr, calD)
	res := float64(resDb)
	// Release Memory
	C.free(unsafe.Pointer(workSpace))
	C.free(unsafe.Pointer(runNo))
	C.free(unsafe.Pointer(prodName))
	C.free(unsafe.Pointer(spCode))
	C.free(unsafe.Pointer(varStr))
	C.free(unsafe.Pointer(calD))

	return res
}

func StoSummaryFloat(entry DllFunc, ws, run, prod, sp, varName, calDate, summaryName string) float64 {
	workSpace := C.CString(ws)
	runNo := C.CString(run)
	prodName := C.CString(prod)
	spCode := C.CString(sp)
	varStr := C.CString(varName)
	calD := C.CString(calDate)
	sName := C.CString(summaryName)

	resDb := C.readStoSummary(entry, workSpace, runNo, prodName, spCode, varStr, calD, sName)
	res := float64(resDb)
	// Release Memory
	C.free(unsafe.Pointer(workSpace))
	C.free(unsafe.Pointer(runNo))
	C.free(unsafe.Pointer(prodName))
	C.free(unsafe.Pointer(spCode))
	C.free(unsafe.Pointer(varStr))
	C.free(unsafe.Pointer(calD))
	C.free(unsafe.Pointer(sName))

	return res
}

/*
type StringStruct struct {
	str unsafe.Pointer
	len int
}

func cChar(s string) *C.char {
	ss := (*StringStruct)(unsafe.Pointer(&s))
	c := (*C.char)(unsafe.Pointer(ss.str))
	return c
}
*/
