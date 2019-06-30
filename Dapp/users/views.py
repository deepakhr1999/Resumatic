from django.shortcuts import render, redirect, get_object_or_404
from django.contrib import messages
from django.contrib.auth.decorators import login_required
from .forms import UserUpdateForm, ProfileUpdateForm, UserRegisterForm
from django.views.generic import ListView, DetailView, CreateView, UpdateView, DeleteView
from .models import *
from django.contrib.auth.models import User
from django.contrib.admin.views.decorators import staff_member_required
from django.utils.decorators import method_decorator
import subprocess
# Create your views here.


def register(request):
    if request.method == "POST":
        form = UserRegisterForm(request.POST)
        if form.is_valid():
            form.save()
            instance = User.objects.get(username=form.cleaned_data["username"])            
            fullname = f"{instance.first_name} {instance.last_name}".replace(" ", "_")
            if form.cleaned_data["is_org"]:# this is and organization
                cmd = ['node', 'node_sdk/invoke.js', 'CreateOrg', instance.username, fullname, instance.password]
                instance.profile.is_org = True
                instance.profile.save()
            else:# normal user
                cmd = ['node', 'node_sdk/invoke.js', 'CreateUser', instance.username, fullname, instance.password, instance.email, "2138973490"]
            cmd = subprocess.run(cmd, encoding="utf-8", stdout=subprocess.PIPE)
            messages.success(request, f'Account has been created!, {cmd.stdout}')

            return redirect('register')
    else:
        form = UserRegisterForm()
    return render(request, 'users/register.html', {'form': form, 'name': 'Register User'})

@login_required
def profile(request):
    if request.method == 'POST':
        u_form = UserUpdateForm(request.POST, instance=request.user)
        p_form = ProfileUpdateForm(request.POST,
                                   request.FILES,
                                   instance=request.user.profile)
        if u_form.is_valid() and p_form.is_valid():
            u_form.save()
            p_form.save()
            messages.success(request, f'Your account has been updated!')
            return redirect('profile')

    else:
        u_form = UserUpdateForm(instance=request.user)
        p_form = ProfileUpdateForm(instance=request.user.profile)

    context = {
        'u_form': u_form,
        'p_form': p_form
    }

    return render(request, 'users/profile.html', context)


