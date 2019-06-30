from django.shortcuts import render, redirect
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.http import HttpResponse
import subprocess
from kyc.models import *
from users.forms import *
import uuid
import json
import time
import requests
# Create your views here.

def welcome(request):
	if(request.user.is_authenticated):
		if request.user.profile.is_org:
			return index_org(request)
		else:
			return index_user(request)
	return render(request, 'kyc/index_user.html', {})

def index_user(request):
	# assume applicant has logged in
	resumes = Resume.objects.filter(user = request.user.username)
	certs = Cert.objects.filter(user = request.user) 	
	return render(request, 'kyc/index_user.html', {"certs": certs, "resumes": resumes})

def index_org(request):
	if request.method == "POST":
		try:
			cert = Cert.objects.get(token = request.POST["token"])
			cert.org_filled = True
			if request.POST["Decision"]=="Confirm":
				cmd = ["node", "node_sdk/invoke.js", "VerifyClaim", cert.token, request.user.username, request.user.password]
				out = subprocess.run(cmd, encoding="utf-8", stdout=subprocess.PIPE).stdout
				messages.success(request, out)
			cert.is_verified = (request.POST["Decision"]=="Confirm")
			cert.save()
		except:
			pass
	
	resumes = Resume.objects.filter(org = request.user.username, sent=True)
	certs = Cert.objects.filter(org = request.user, is_verified=False, org_filled= False)
	return render(request, 'kyc/index_org.html', {"certs": certs, "resumes": resumes})

def resume(request, resume_id=None):
	if request.method == "POST":
		if resume_id==None:
			# create resume
			tokens = request.POST.getlist("tokens")
			tokens_str = json.dumps(tokens)
			resume = Resume(user=request.user.username, tokens=tokens_str)
			resume.save()
			resume_id = resume.id
		else:
			#update resume
			resume = Resume.objects.get(id=resume_id)
			if "org" in request.POST:
				resume.org = request.POST["org"]
			elif "sent" in request.POST:
				resume.sent = True
				messages.success(request, f"Resume has been sent successfully to {resume.org}")
			elif "delete" in request.POST:
				resume.delete()
				return redirect("welcome")
			resume.save()
	if resume_id!=None:
		resume = Resume.objects.get(id=resume_id).__dict__
		resume["certs"] = Cert.objects.filter(token__in=json.loads(resume["tokens"]))
		resume["user"] = User.objects.get(username=resume["user"])
		return render(request, 'kyc/resume.html', resume)
	return redirect("welcome")


@login_required
def make_claim(request):
	if request.method == "POST":
		if request.user.profile.is_org:
			form = ProvideCertForm(request.POST)
		else:
			form = RegisterCertForm(request.POST)
		if form.is_valid():
			cert = form.save(commit=False)
			if request.user.profile.is_org:
				cert.org = request.user
			else:
				cert.user = request.user
			cert.token = str(uuid.uuid4())
			cert.save()
			#push it to blockchain
			cmd = ["node", "node_sdk/invoke.js", "MakeClaim", cert.token, cert.user.username, cert.user.password, cert.org.username, cert.title, str(cert.time_till)]
			print(" ".join(cmd))
			out = subprocess.run(cmd, encoding="utf-8", stdout=subprocess.PIPE).stdout
			messages.success(request, out)
			return redirect('welcome')
	else:
		form = RegisterCertForm()
	return render(request, 'users/register.html', {'form': form, 'name': 'Add Credentials'})

def exhaustiveQuery(certs):
	context = []
	# print("Hello=================================")
	for cert in certs:
		cmd = ['node', 'node_sdk/query.js', 'Query', cert.token]
		out = subprocess.run(cmd ,encoding="utf-8", stdout=subprocess.PIPE)
		out = out.stdout
		out = out[out.find("{"):]
		context.append( (cert.token, json.loads(out)) )
		
	return context

def Query(request, token):
	cmd = ['node', 'node_sdk/query.js', 'Query', token]
	out = subprocess.run(cmd ,encoding="utf-8", stdout=subprocess.PIPE)
	out = out.stdout
	out = json.loads(out[out.find("{"):])
	return HttpResponse(out["IsVerified"])
